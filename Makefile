BINARY=go-microservice

DOCKER_IMAGE_NAME=greg-nicolle/go-microservice

.DEFAULT_GOAL: ${BINARY}

${BINARY}:
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BINARY} .

buildandrun:
	@make clean && go build -a -installsuffix cgo -o ${BINARY} . && ./${BINARY}

check: test lint vet

test:
	@go test -race -v $(shell go list ./... | grep -v /vendor/)

lint:
	@go list ./...  | grep -v /vendor/ |  xargs -L1 golint -set_exit_status

vet:
	@go vet $(shell go list ./... | grep -v /vendor/)

install:
	@go install $(shell go list ./... | grep -v /vendor/)

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker-build: ${BINARY}
	@docker build -t ${DOCKER_IMAGE_NAME} .