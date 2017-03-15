# go-microservice
[![Build Status](https://travis-ci.org/greg-nicolle/go-microservice.svg?branch=master)](https://travis-ci.org/greg-nicolle/go-microservice)
## Build
```shell
make
```

## Test app
```shell
make check
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

##