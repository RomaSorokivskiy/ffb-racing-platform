package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"
  "time"

  "ffb.local/matchmaker/internal/rooms"
)

type claimReq struct {
  UserID string `json:"userId"`
  TTLsec int64  `json:"ttlSec"` // optional
}
type releaseReq struct {
  UserID string `json:"userId"`
  CarID  string `json:"carId"`
}
type hbReq struct {
  CarID string `json:"carId"`
  Busy  bool   `json:"busy"`
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

func sseHeaders(w http.ResponseWriter) (http.Flusher, bool) {
  w.Header().Set("Content-Type", "text/event-stream")
  w.Header().Set("Cache-Control", "no-cache")
  w.Header().Set("Connection", "keep-alive")
  flusher, ok := w.(http.Flusher)
  return flusher, ok
}

func main() {
  reg := rooms.NewRegistry(5)
  defer reg.Close()

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
    ttl := time.Duration(req.TTLsec) * time.Second
    car, err := reg.Claim(req.UserID, ttl)
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

  mux.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
      w.WriteHeader(http.StatusMethodNotAllowed)
      return
    }
    var req hbReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CarID == "" {
      w.WriteHeader(http.StatusBadRequest)
      return
    }
    var err error
    if req.Busy {
      err = reg.MarkBusy(req.CarID)
    } else {
      err = reg.MarkFree(req.CarID)
    }
    if err != nil {
      http.Error(w, err.Error(), http.StatusNotFound)
      return
    }
    w.WriteHeader(http.StatusNoContent)
  })

  // SSE stream
  mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
    flusher, ok := sseHeaders(w)
    if !ok {
      http.Error(w, "stream unsupported", http.StatusInternalServerError)
      return
    }

    ch := reg.Subscribe()
    defer reg.Unsubscribe(ch)

    notify := r.Context().Done()
    for {
      select {
      case <-notify:
        return
      case ev, ok := <-ch:
        if !ok {
          return
        }
        fmt.Fprintf(w, "data: %s\n\n", rooms.MarshalEvent(ev))
        flusher.Flush()
      }
    }
  })

  addr := ":8081"
  if v := os.Getenv("PORT"); v != "" {
    addr = ":" + v
  }
  log.Println("matchmaker listening", addr)
  log.Fatal(http.ListenAndServe(addr, cors(mux)))
}
