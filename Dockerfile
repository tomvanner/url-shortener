FROM golang:1.23 AS build

# go-sqlite3 is a CGO enabled package
ENV CGO_ENABLED=1

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o app main.go

EXPOSE 8000

ENTRYPOINT ["./app"]