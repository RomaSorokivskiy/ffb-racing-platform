package main

import (
	"log"
	"net/http"

	"ffb.local/gateway/internal/signaling"
)

func main() {
	s := signaling.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.Health)
	mux.HandleFunc("/signal", s.Signal)

	handler := s.CORS(s.Log(mux))
	srv := &http.Server{Addr: ":8080", Handler: handler}

	log.Println("gateway listening :8080")
	log.Fatal(srv.ListenAndServe())
}
