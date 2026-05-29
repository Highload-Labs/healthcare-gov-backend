FROM golang:1.26.3-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o healthcare-gov cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/healthcare-gov .

EXPOSE 8080

CMD ["./healthcare-gov"]