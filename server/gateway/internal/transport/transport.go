package transport

// NOTE: This file is a stub to pin the dependency direction and
// to give a clear place to implement WebRTC session creation.
// In the next step, we will wire pion/webrtc here and expose an API
// like: NewPeerFor(carID string) (*Peer, error), Peer.BindTracks(...)

type Peer struct {
	CarID string
	// pc *webrtc.PeerConnection
}

func NewPeer(carID string) (*Peer, error) {
	return &Peer{CarID: carID}, nil
}

func (p *Peer) Close() error {
	// if p.pc != nil { return p.pc.Close() }
	return nil
}
