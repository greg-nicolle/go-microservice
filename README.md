# go-microservice

## Build
```shell
go build *.go
```

## Run

To run an instance do
```shell
go run *.go -listen=:8080 &
```

to run multiple instance do
```shell
go run *.go -listen=:8001 &
go run *.go -listen=:8001 &
go run *.go -listen=:8080 -proxy=localhost:8001,localhost:800 &
```