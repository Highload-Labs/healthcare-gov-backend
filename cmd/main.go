package main

import (
	"github.com/Highload-Labs/healthcare-gov-backend/internal/transport"
)

func main() {
	server := transport.NewHTTP()
	server.SetupAndServe()
}
