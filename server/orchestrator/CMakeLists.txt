cmake_minimum_required(VERSION 3.20)
project(orchestrator CXX)
set(CMAKE_CXX_STANDARD 20)


add_executable(orchestrator
src/main.cpp
src/ffb_model.cpp
)


# TODO: link protobuf/grpc if needed
