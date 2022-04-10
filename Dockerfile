from golang:1.17.3-alpine3.14

run mkdir /app
workdir /app
run go get -u github.com/gorilla/mux
#postgres driver
run go get -u github.com/lib/pq
#golang query builder
run go get -u github.com/doug-martin/goqu/v9
#hash password
run go get -u golang.org/x/crypto/bcrypt
copy . .
run go build app.go
cmd ["./app"]