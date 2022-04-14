FROM golang:latest

LABEL maintainer = "Pavel Erokgin <pavel.v.erokhin@gmail.com>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build

CMD ["./user-microservice-go"]
