# Golang API Login-Logout-Transaction

### Instalasi
Repo ini menggunakan docker. Jadi untuk menjalankan program hanya memerlukan command
```bash
docker-compose up -d
```

### Contoh request
- Login `GET` ke alamat `http://localhost:8000/api/login`  
Pada postman, masukkan ke dalam `Body > raw > JSON`  
    ```json
    {
        "username": "hero",
        "password": "hero"
    }
    ```
    response  
    ```json
    {
        "type": "login",
        "message": "sukses"
    }
    ```
    jika login berlebihan, maka muncul response  
    ```json
    {
        "type": "login",
        "message": "user already login"
    }
    ```
- Logout `GET` ke alamat `http://localhost:8000/api/logout`  
Pada postman, masukkan ke dalam `Body > raw > JSON`
    ```json
    {
        "username": "hero",
        "password": "hero"
    }
    ```
    response  
    ```json
    {
        "type": "logout",
        "message": "sukses"
    }
    ```
    jika logout berlebihan, maka muncul response  
    ```json
    {
        "type": "logout",
        "message": "user already logout"
    }
    ```
- Payment `POST` ke alamat `http://localhost:8000/api/payment`  
pada postman, masukkan ke dalam `Body > raw > JSON`
    ```json
    {
        "from": "hero",
        "to": "zeus",
        "amount": 10000
    }
    ```
    pastikan untuk login terlebih dulu, lalu lanjutkan untuk mulai melakukan transfer.  
    response jika sudah login
    ```json
    {
        "type": "payment",
        "message": "transaction success"
    }
    ```
    response jika sebelum login
    ```json
    {
        "type": "payment",
        "message": "user need to login"
    }
    ```