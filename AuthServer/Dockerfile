FROM golang:1.20

WORKDIR /usr/src/app

COPY main.go main.go
COPY users.db users.db
RUN go mod init main && go mod tidy
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]