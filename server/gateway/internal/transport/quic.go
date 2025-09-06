package transport

// This is a stub for QUIC control/FFB channel using quic-go.
// We'll fill it with real dial/accept logic later.

type QUICEndpoint struct {
	Addr string
}

func NewQUIC(addr string) *QUICEndpoint {
	return &QUICEndpoint{Addr: addr}
}

func (q *QUICEndpoint) Start() error {
	// TODO: start listener for control/ffb streams
	return nil
}

func (q *QUICEndpoint) Stop() error {
	// TODO: gracefully stop
	return nil
}
