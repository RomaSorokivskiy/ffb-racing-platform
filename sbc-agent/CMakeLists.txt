cmake_minimum_required(VERSION 3.20)
project(sbc_agent CXX)
set(CMAKE_CXX_STANDARD 20)


add_executable(sbc_agent
src/main.cpp
src/capture.cpp
src/webrtc_out.cpp
src/telemetry_bridge.cpp
)
