FROM scratch

ADD go-microservice /
ADD config.yml /

EXPOSE 8080

CMD ["/go-microservice", "--configPath", "/config.yml", "--service", "all"]