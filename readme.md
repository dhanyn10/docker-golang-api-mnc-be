# GOLANG API

### Petunjuk
Berikut ini list dari request yang bisa digunakan beserta contoh hasilnya (dalam postman):
- post new article  
    contoh cara eksekusi
    ```
    http://localhost:8000/api/article?author=jostein&title=dunia suphie&body=novel tentang dunia filsafat untuk anak-anak
    ```
- get article dengan sorting berdasar parameter nama kolom `param` dan nilainya `query`
    ```
    http://localhost:8000/api/article?param=author&query=suparman
    ```
    
    hasil
    ```
    [
        {
            "id": 20,
            "author": "suparman",
            "title": "perjuangan",
            "body": "perang dunia dua menuju akhir",
            "created": "2021-11-28T15:59:09.658958Z"
        }
    ]
    ```
- get list all articles dengan sorting berdasarkan waktu pembuatan terbaru
    ```
    http://localhost:8000/api/articles
    ```
    hasil
    ```
    [
        {
            "id": 21,
            "author": "jostein",
            "title": "dunia suphie",
            "body": "novel tentang dunia filsafat untuk anak-anak",
            "created": "2021-11-28T16:13:20.703631Z"
        },
        {
            "id": 20,
            "author": "suparman",
            "title": "perjuangan",
            "body": "perang dunia dua menuju akhir",
            "created": "2021-11-28T15:59:09.658958Z"
        }
    ]
    ```
