FROM scratch
ADD go-microservice /
EXPOSE 8080
CMD ["/go-microservice"]