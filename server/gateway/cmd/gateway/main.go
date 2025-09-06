package main


import (
"log"
"net/http"
)


func main() {
mux := http.NewServeMux()
// TODO: add /signal (WebRTC) and /quic endpoints
mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
_, _ = w.Write([]byte("ok"))
})


srv := &http.Server{Addr: ":8080", Handler: mux}
log.Println("gateway listening :8080")
log.Fatal(srv.ListenAndServe())
}
