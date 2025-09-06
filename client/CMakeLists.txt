cmake_minimum_required(VERSION 3.20)
project(client CXX)
set(CMAKE_CXX_STANDARD 20)


add_executable(client
src/main.cpp
src/video_player.cpp
src/ffb_driver.cpp
src/net_client.cpp
)
