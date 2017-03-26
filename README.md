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

To run all service on one instance
```shell
go run *.go --configPath ./config.yml --service all
```

to run multiple instance do
```shell
go run *.go --configPath ./config.yml --service master
go run *.go --configPath ./config.yml --service string
go run *.go --configPath ./config.yml --service wikidata
```

##