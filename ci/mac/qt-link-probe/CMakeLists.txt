cmake_minimum_required(VERSION 3.21)
project(qt_link_probe LANGUAGES CXX)

find_package(Qt6 REQUIRED COMPONENTS Widgets Network Multimedia MultimediaWidgets)

qt_add_executable(qt_link_probe main.cpp)
target_link_libraries(
  qt_link_probe
  PRIVATE Qt6::Widgets Qt6::Network Qt6::MultimediaWidgets
)
