package signaling

import (
	"encoding/json"
	"log"
	"net/http"

	"ffb.local/gateway/internal/transport"
)

type Offer struct {
	SDP string `json:"sdp"`
}
type Answer struct {
	SDP string `json:"sdp"`
}

type Handler struct {
	wrtc *transport.WebRTC
}

func New() *Handler {
	return &Handler{wrtc: transport.NewWebRTC()}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) Signal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var off Offer
	if err := json.NewDecoder(r.Body).Decode(&off); err != nil || off.SDP == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ansSDP, err := h.wrtc.HandleOffer(r.Context(), off.SDP)
	if err != nil {
		log.Println("webrtc error:", err)
		http.Error(w, "webrtc failed", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(Answer{SDP: ansSDP})
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h *Handler) CORS(next http.Handler) http.Handler {
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

func (h *Handler) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
