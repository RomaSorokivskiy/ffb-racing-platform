# FFB Racing Platform — 5-Car MVP

Real-time teleoperation racing system: physical RC cars (with camera/IMU), low-latency video via WebRTC, control & Force Feedback via QUIC/UDP, desktop client with HUD and FFB, launcher, and web portal.

## Goals (MVP)

- ≤60 ms P95 camera→screen latency on LAN
- Stable FFB (250–1000 Hz) without oscillations
- 5 cars concurrently: lobby + per-room orchestration

## Tech

- **Proto**: protobuf (shared contracts)
- **Server**: Go (gateway, matchmaker), C++ (orchestrator)
- **Client**: C++ (video, HUD, FFB)
- **SBC Agent**: C++ (capture/encode, telemetry)
- **Launcher**: Electron + TS
- **Web**: Next.js + TS

## Quick start

```bash
# 1) Generate protobuf stubs (Go/C++)
make proto


# 2) Build C++ targets (client, orchestrator, sbc-agent)
make build-cpp


# 3) Build Go services (gateway, matchmaker)
make build-go


# 4) Run docker-compose for infra (dev)
make up
```
