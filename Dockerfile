FROM golang:1.26.3-alpine

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o healthcare-gov cmd/main.go

EXPOSE 8080

CMD ["./healthcare-gov"]