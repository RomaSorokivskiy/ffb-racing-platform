package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"ffb.local/matchmaker/internal/rooms"
)

type claimReq struct {
	UserID string `json:"userId"`
}
type releaseReq struct {
	UserID string `json:"userId"`
	CarID  string `json:"carId"`
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("access-control-allow-origin", "*")
		w.Header().Set("access-control-allow-headers", "content-type")
		w.Header().Set("access-control-allow-methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	reg := rooms.NewRegistry(5)

	mux := http.NewServeMux()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(reg.List())
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	})

	mux.HandleFunc("/claim", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req claimReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		car, err := reg.Claim(req.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		resp, _ := json.Marshal(car)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	})

	mux.HandleFunc("/release", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req releaseReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.CarID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		car, err := reg.Release(req.UserID, req.CarID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		resp, _ := json.Marshal(car)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	})

	addr := ":8081"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}
	log.Println("matchmaker listening", addr)
	log.Fatal(http.ListenAndServe(addr, cors(mux)))
}
