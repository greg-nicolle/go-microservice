FROM scratch
ADD go-microservice /
EXPOSE 8080
CMD ["/go-microservice -listen=:8001","/go-microservice -listen=:8002","/go-microservice -listen=:8080 -proxy=localhost:8001,localhost:800"]