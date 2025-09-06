# ADR-0001: Transport (WebRTC + QUIC)

**Decision**: Use WebRTC (video SRTP) and QUIC/UDP (control + FFB) with protobuf payloads.

**Rationale**: low-latency, packet recovery (NACK/PLI), prioritization for FFB.

**Consequences**: need TURN/STUN infra; QUIC libs integration on SBC & client.
