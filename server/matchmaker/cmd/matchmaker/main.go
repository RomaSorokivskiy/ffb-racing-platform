package main


import (
"log"
"net/http"
)


func main() {
mux := http.NewServeMux()
mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
// TODO: list/create rooms for 5 cars
w.WriteHeader(http.StatusOK)
_, _ = w.Write([]byte("[]"))
})
log.Println("matchmaker listening :8081")
log.Fatal(http.ListenAndServe(":8081", mux))
}
