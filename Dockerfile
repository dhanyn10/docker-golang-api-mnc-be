from golang:1.17.3-alpine3.14

run mkdir /app
workdir /app
run go get -u github.com/gorilla/mux
run go get -u github.com/lib/pq
#golang query builder
run go get -u github.com/doug-martin/goqu/v9
copy . .
run go build app.go
cmd ["./app"]