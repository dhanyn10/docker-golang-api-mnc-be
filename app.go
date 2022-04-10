package main

import (
	"encoding/json"
    "encoding/hex"
	"math/rand"
	"log"
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/doug-martin/goqu/v9"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

const (
    host     = "db"
    port     = 5432
    user     = "admin"
    password = "admin"
    dbname   = "admin"
)

// Nasabah Struct(Model)
type Nasabah struct {
	Id			int  	`json:"id"`
	Username	string	`json:"username"`
	Password	string	`json:"password"`
	Token		string 	`json:"token"`
	Tabungan	int 	`json:"tabungan"`
}

// Transaksi Struct(Model)
type Transaksi struct {
	Id			int  	`json:"id"`
	From		string	`json:"from"`
	To			string	`json:"to"`
	Amount		int 	`json:"amount"`
	Date		string 	`json:"date"`
}
// Report Struct(Model)
type Report struct {
	DataType 	string `json:"type"`
	Message		string `json:"message"`
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func connectDB() *sql.DB{
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	return db
}
/**
 * transaksi	: menunjukkan event yang sedang dilakukan sesuai dengan nama fungsinya
 * 				-> login | logout | payment
 * cond			: menunjukkan kondisi keberhasilan transaksi setiap menekan <Send>
 * 				-> 1 (sukses) | 0 (gagal)
 * activity		: menunjukkan activity yang sedang berlangsung, baik itu activity
 * 				  yang dijalankan terhadap database (crud) maupun error message yang diterima
 * 				  saat menjalankan transaksi.
 */
func ActivityHistory(transaksi string, cond int, activity string) {
	db := connectDB()
	//history
	historyData, _, _:= goqu.Insert("history").
	Cols("event", "cond", "activity", "datetime").
	Vals(goqu.Vals{transaksi, cond, activity, time.Now()}).ToSQL()
	_, errQuery := db.Query(historyData)
	CheckError(errQuery)
}

//Login
func Login(w http.ResponseWriter, r *http.Request) {

	LoginReport := Report{}
	rmsg := ""

	w.Header().Set("Content-Type", "application/json")

	reqData, err := ioutil.ReadAll(r.Body)
	if err != nil{
		log.Fatal(err)
		return
	}
	var nasabah Nasabah
	// unmarshal
	if err := json.Unmarshal(reqData, &nasabah); err != nil {
		panic(err)
	}
	CheckError(err)

	db := connectDB()
	selectUser, _, _ := goqu.From("nasabah").Where(goqu.Ex{"username": nasabah.Username}).ToSQL()
	rows, err := db.Query(selectUser)
	CheckError(err)
	defer rows.Close()

	registered := false
	loginToken := ""
	loginPassword := ""
	for rows.Next() {
		var id int
		var username string
		var password string
		var token string
		var tabungan int
	
		err = rows.Scan(&id, &username, &password, &token, &tabungan)
		CheckError(err)
		registered = true
		loginToken = token
		loginPassword = password
	}

	//lolos verifikasi cek username 
	if registered == true {

		//password benar
		if CheckPasswordHash(nasabah.Password, loginPassword) == true {	
			//token masing kosong, user baru pertama kali login
			if loginToken == "" {
				goQuery, _, _ := goqu.Update("nasabah").Set(goqu.Record{"token": GenerateSecureToken(5)}).Where(goqu.Ex{"username": nasabah.Username}).ToSQL()
				_, err := db.Query(goQuery)
				CheckError(err)
				rmsg = "success"
				//history
				ActivityHistory("login", 1, goQuery)
			} else {
				rmsg = "user already login"
				//history
				ActivityHistory("logout", 0, rmsg)
			}
		} else {
			//password tidak sama
			rmsg = "wrong password"
			//history
			ActivityHistory("logout", 0, rmsg)
		}
	} else {
		rmsg = "user not registered"
		//history
		ActivityHistory("logout", 0, rmsg)
	}

	LoginReport.DataType = "login"
	LoginReport.Message = rmsg
	reportJson, reportErr := json.Marshal(LoginReport)
	if reportErr != nil {
		CheckError(reportErr)
	}
	w.Write(reportJson)
}

//Logout
func Logout(w http.ResponseWriter, r *http.Request) {

	LogoutReport := Report{}
	rmsg := ""

	w.Header().Set("Content-Type", "application/json")

	reqData, err := ioutil.ReadAll(r.Body)
	if err != nil{
		log.Fatal(err)
		return
	}
	var nasabah Nasabah
	// unmarshal
	if err := json.Unmarshal(reqData, &nasabah); err != nil {
		panic(err)
	}
	CheckError(err)

	db := connectDB()
	selectUser, _, _ := goqu.From("nasabah").Where(goqu.Ex{"username": nasabah.Username}).ToSQL()
	rows, err := db.Query(selectUser)
	CheckError(err)
	defer rows.Close()

	registered := false
	logoutPassword := ""
	logoutToken := ""
	for rows.Next() {
		var id int
		var username string
		var password string
		var token string
		var tabungan int

	
		err = rows.Scan(&id, &username, &password, &token, &tabungan)
		CheckError(err)
		logoutPassword = password
		logoutToken = token
		registered = true
	}

	//username terdaftar
	if registered == true {
		//password benar
		if CheckPasswordHash(nasabah.Password, logoutPassword) {
			//cek token
			if logoutToken != "" {
				goQuery, _, _ := goqu.Update("nasabah").Set(goqu.Record{"token": ""}).Where(goqu.Ex{"username": nasabah.Username}).ToSQL()
				_, err := db.Query(goQuery)
				CheckError(err)
				rmsg = "success"
				//history
				ActivityHistory("logout", 1, goQuery)
			} else {
				rmsg = "user already logout"
				//history
				ActivityHistory("logout", 0, rmsg)
			}
		} else {
			rmsg = "wrong password"
			//history
			ActivityHistory("logout", 0, rmsg)
		}
	} else {
		rmsg = "user not registered"
		//history
		ActivityHistory("logout", 0, rmsg)
	}
	
	LogoutReport.DataType = "logout"
	LogoutReport.Message = rmsg
	reportJson, reportErr := json.Marshal(LogoutReport)
	if reportErr != nil {
		CheckError(reportErr)
	}
	w.Write(reportJson)
}

//Payment
func Payment(w http.ResponseWriter, r *http.Request) {

	PaymentReport := Report{}
	rmsg := ""

	w.Header().Set("Content-Type", "application/json")

	reqData, err := ioutil.ReadAll(r.Body)
	if err != nil{
		log.Fatal(err)
		return
	}
	var transaksi Transaksi
	// unmarshal
	if err := json.Unmarshal(reqData, &transaksi); err != nil {
		panic(err)
	}
	CheckError(err)

	db := connectDB()
	datauser, _, _ := goqu.From("nasabah").Where(goqu.Ex{"username": transaksi.From}).ToSQL()
	//cek user sudah login
	rowsFrom, err := db.Query(datauser)
	defer rowsFrom.Close()

	registered:= false
	userToken := ""
	tabunganFrom := 0
	for rowsFrom.Next() {
		var id int
		var username string
		var password string
		var token string
		var tabungan int

		//data user from
		err = rowsFrom.Scan(&id, &username, &password, &token, &tabungan)
		CheckError(err)
		userToken = token
		tabunganFrom = tabungan
		registered = true
	}
	
		//jika token tidak kosong, maka user telah login
		//bisa dilanjutkan ke langkah berikutnya
		if userToken != "" {
			//membutuhkan data user to untuk
			//memastikan sudah terdaftar di database
			uToExist := false

			userTo, _, _ := goqu.From("nasabah").Where(goqu.Ex{"username": transaksi.To}).ToSQL()
			rowsTo, err := db.Query(userTo)
			CheckError(err)
			defer rowsTo.Close()

			for rowsTo.Next() {
				uToExist = true
			}

			//user to terdaftar di database
			if uToExist == true {
				//masukkan data transaksi

				//validasi nilai transfer tidak lebih dari nilai saldo userFrom
				if transaksi.Amount <= tabunganFrom {
					//lolos, nilai transaksi aman
					transaksiQuery, _, _:= goqu.Insert("transaksi").
					Cols("from", "to", "amount", "datetime").
					Vals(goqu.Vals{transaksi.From, transaksi.To, transaksi.Amount, time.Now()}).ToSQL()
					_, errQuery := db.Query(transaksiQuery)
					CheckError(errQuery)

					//history
					ActivityHistory("payment", 1, transaksiQuery)

					dataTabunganTo, _, _ := goqu.From("nasabah").Where(goqu.Ex{"username": transaksi.To}).ToSQL()
					//cek user sudah login
					rowsTabunganTo, errDataTabunganTo := db.Query(dataTabunganTo)
					CheckError(errDataTabunganTo)
					defer rowsTabunganTo.Close()

					tabunganTo := 0
					for rowsTabunganTo.Next() {
						var id int
						var username string
						var password string
						var token string
						var tabungan int
				
						//ambil data tabungan dari database
						err = rowsTabunganTo.Scan(&id, &username, &password, &token, &tabungan)
						//atur nilai awal dengan nilai di database
						tabunganTo = tabungan
					}
					//userFrom: kurangi nilai tabungan dengan nilai transaksi
					tabunganFrom = tabunganFrom - transaksi.Amount
					//userTo: tambah nilai tabungan dengan nilai transaksi
					tabunganTo = tabunganTo + transaksi.Amount

					updateTabunganFrom, _, _ := goqu.Update("nasabah").
					Set(goqu.Record{"tabungan": tabunganFrom}).
					Where(goqu.Ex{"username": transaksi.From}).
					ToSQL()
					_, errTabunganFrom := db.Query(updateTabunganFrom)
					CheckError(errTabunganFrom)
					//history
					ActivityHistory("payment", 1, updateTabunganFrom)

					updateTabunganTo, _, _ := goqu.Update("nasabah").
					Set(goqu.Record{"tabungan": tabunganTo}).
					Where(goqu.Ex{"username": transaksi.To}).
					ToSQL()
					_, errTabunganTo := db.Query(updateTabunganTo)
					CheckError(errTabunganTo)
					
					//history
					ActivityHistory("payment", 1, updateTabunganTo)
					rmsg = "transaction success"
				} else {
					rmsg = "you dont have enough money"
					//history
					ActivityHistory("payment", 0, rmsg)
				}
			} else {
				rmsg = "canceled, data not valid"
				//history
				ActivityHistory("payment", 0, rmsg)
			}
		} else {
			rmsg = "user need to login"
			//history
			ActivityHistory("payment", 0, rmsg)
		}

	if registered == false {
		rmsg = "user not registered"
		//history
		ActivityHistory("payment", 1, rmsg)
	}
	PaymentReport.DataType = "payment"
	PaymentReport.Message = rmsg
	reportJson, reportErr := json.Marshal(PaymentReport)
	if reportErr != nil {
		CheckError(reportErr)
	}
	w.Write(reportJson)
}

func CheckError(err error) {
    if err != nil {
        fmt.Println(err)
    }
}

func GenerateSecureToken(length int) string {
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        return ""
    }
    return hex.EncodeToString(b)
}

func main() {
	//init router
	r:= mux.NewRouter()

	r.HandleFunc("/api/login", Login).Methods("GET")
	r.HandleFunc("/api/logout", Logout).Methods("GET")
	r.HandleFunc("/api/payment", Payment).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
