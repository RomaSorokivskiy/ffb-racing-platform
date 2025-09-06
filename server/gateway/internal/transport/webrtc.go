package transport

import (
	"context"

	"github.com/pion/webrtc/v3"
)

type WebRTC struct{}

func NewWebRTC() *WebRTC { return &WebRTC{} }

// HandleOffer: приймає SDP offer, повертає SDP answer.
// Поки без медіа — чистий SDP обмін (перевірка траси).
func (w *WebRTC) HandleOffer(ctx context.Context, offerSDP string) (string, error) {
	m := webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		return "", err
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&m))
	pc, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return "", err
	}
	defer pc.Close()

	// (опціонально) заглушки обробників
	pc.OnICEConnectionStateChange(func(s webrtc.ICEConnectionState) { _ = s })

	offer := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offerSDP}
	if err := pc.SetRemoteDescription(offer); err != nil {
		return "", err
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return "", err
	}
	if err := pc.SetLocalDescription(answer); err != nil {
		return "", err
	}

	// Повертаємо SDP answer
	return answer.SDP, nil
}
