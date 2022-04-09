package main

import (
	"encoding/json"
    "encoding/hex"
	"math/rand"
	"log"
	"time"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
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

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)

	selectUser, _, _ := goqu.From("nasabah").Where(goqu.Ex{"username": nasabah.Username, "password": nasabah.Password}).ToSQL()
	rows, err := db.Query(selectUser)
	CheckError(err)
	registered := false

	defer rows.Close()
	for rows.Next() {
		var id int
		var username string
		var password string
		var token string
		var tabungan int

	
		err = rows.Scan(&id, &username, &password, &token, &tabungan)
		CheckError(err)
		if token == "" {
			goQuery, _, _ := goqu.Update("nasabah").Set(goqu.Record{"token": GenerateSecureToken(5)}).Where(goqu.Ex{"username": nasabah.Username}).ToSQL()
			_, err := db.Query(goQuery)
			CheckError(err)
			rmsg = "sukses"
		} else {
			rmsg = "user already login"
		}
		registered = true
	}

	if registered == false {
		rmsg = "user not registered"
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

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)

	selectUser := fmt.Sprintf("SELECT * FROM nasabah WHERE username= '%s' AND password='%s'", nasabah.Username, nasabah.Password)
	rows, err := db.Query(selectUser)
	CheckError(err)
	defer rows.Close()

	registered := false
	for rows.Next() {
		var id int
		var username string
		var password string
		var token string
		var tabungan int

	
		err = rows.Scan(&id, &username, &password, &token, &tabungan)
		CheckError(err)
		if token != "" {
			goQuery, _, _ := goqu.Update("nasabah").Set(goqu.Record{"token": ""}).Where(goqu.Ex{"username": nasabah.Username}).ToSQL()
			_, err := db.Query(goQuery)
			CheckError(err)
			rmsg = "sukses"
		} else {
			rmsg = "user already logout"
		}
		registered = true
	}

	if registered == false {
		rmsg = "user not registered"
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

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)

	datauser, _, _ := goqu.From("nasabah").Where(goqu.Ex{"username": transaksi.From}).ToSQL()
	//cek user sudah login
	rows, err := db.Query(datauser)
	defer rows.Close()

	registered:= false
	for rows.Next() {
		var id int
		var username string
		var password string
		var token string
		var tabungan int

		//data user from
		err = rows.Scan(&id, &username, &password, &token, &tabungan)
		CheckError(err)

		//jika token tidak kosong, maka user telah login
		//bisa dilanjutkan ke langkah berikutnya
		if token != "" {
			//membutuhkan data user from dan to untuk
			//memastikan sudah terdaftar di database
			uFromExist := false
			uToExist := false

			userFrom, _, _ := goqu.From("nasabah").Where(goqu.Ex{"username": transaksi.From}).ToSQL()
			rowsFrom, err := db.Query(userFrom)
			CheckError(err)
			defer rowsFrom.Close()

			for rowsFrom.Next() {
				uFromExist = true
			}

			userTo, _, _ := goqu.From("nasabah").Where(goqu.Ex{"username": transaksi.To}).ToSQL()
			rowsTo, err := db.Query(userTo)
			CheckError(err)
			defer rowsTo.Close()

			for rowsTo.Next() {
				uToExist = true
			}

			//user from dan user to terdaftar di database
			if uFromExist == true && uToExist == true {
				//masukkan data transaksi

				//validasi nilai transfer tidak lebih dari nilai saldo userFrom
				if transaksi.Amount <= tabungan {
					//lolos, nilai transaksi aman
					transaksiQuery, _, _:= goqu.Insert("transaksi").
					Cols("from", "to", "amount", "datetime").
					Vals(goqu.Vals{transaksi.From, transaksi.To, transaksi.Amount, time.Now()}).ToSQL()
					_, errQuery := db.Query(transaksiQuery)
					CheckError(errQuery)
					
					//userFrom: kurangi nilai tabungan dengan nilai transaksi
					tabunganFrom := tabungan - transaksi.Amount
					//userTo: tambah nilai tabungan dengan nilai transaksi
					tabunganTo := tabungan + transaksi.Amount

					updateTabunganFrom, _, _ := goqu.Update("nasabah").Set(goqu.Record{"tabungan": tabunganFrom}).Where(goqu.Ex{"username": transaksi.From}).ToSQL()
					_, errTabunganFrom := db.Query(updateTabunganFrom)
					CheckError(errTabunganFrom)

					updateTabunganTo, _, _ := goqu.Update("nasabah").Set(goqu.Record{"tabungan": tabunganTo}).Where(goqu.Ex{"username": transaksi.To}).ToSQL()
					_, errTabunganTo := db.Query(updateTabunganTo)
					CheckError(errTabunganTo)
					rmsg = "transaction success"
				} else {
					rmsg = "you dont have enough money"
				}
			} else {
				rmsg = "canceled, data not valid"
			}
		} else {
			rmsg = "user need to login"
		}
		registered = true
	}

	if registered == false {
		rmsg = "user not registered"
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
