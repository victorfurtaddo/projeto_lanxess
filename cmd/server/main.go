package main

import (
	"log"
	"net/http"
	"os"

	"projeto-lanxess/internal/api"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := api.NewServer()
	log.Printf("servidor iniciado em http://localhost%s", addr)
	if err := http.ListenAndServe(addr, server.Routes()); err != nil {
		log.Fatal(err)
	}
}
