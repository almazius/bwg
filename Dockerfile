FROM golang:latest

WORKDIR /bwg2/

COPY .. .
COPY go.mod .
COPY go.sum .
RUN go mod download
EXPOSE 8080

RUN go build -o main cmd/main.go

CMD ["./main"]