FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ENV CONFIG_DIRECTORY=./config

COPY . .

RUN go test ./...

RUN go build -o app main.go

CMD ["./app"]