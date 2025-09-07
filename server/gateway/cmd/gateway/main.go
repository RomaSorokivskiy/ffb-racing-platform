package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"ffb.local/gateway/internal/auth"
	"ffb.local/gateway/internal/signaling"
)

type sessionReq struct {
	UserID string `json:"userId"`
	CarID  string `json:"carId"`
}

func main() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}
	signer := auth.NewSigner(secret, 10*time.Minute)

	s := signaling.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.Health)
	mux.HandleFunc("/signal", s.Signal)

	// Create session token (bind user + car) -> returns JWT
	mux.HandleFunc("/session/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req sessionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.CarID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tok, err := signer.Sign(req.UserID, req.CarID)
		if err != nil {
			http.Error(w, "token error", http.StatusInternalServerError)
			return
		}
		resp := map[string]string{"token": tok}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	handler := s.CORS(s.Log(mux))
	srv := &http.Server{Addr: ":8080", Handler: handler}

	log.Println("gateway listening :8080")
	log.Fatal(srv.ListenAndServe())
}
