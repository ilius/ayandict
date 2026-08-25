cmake_minimum_required(VERSION 3.21)
project(qt_link_probe LANGUAGES CXX)

if(NOT DEFINED AYANDICT_QT_VERSION)
  message(FATAL_ERROR "AYANDICT_QT_VERSION must be set")
endif()

# Windows-GNU defaults to import libraries before static archives. Override it
# after project() loads that platform file and before any dependency lookup.
set(CMAKE_FIND_LIBRARY_SUFFIXES ".a")

find_package(PkgConfig REQUIRED)
pkg_check_modules(
  AYANDICT_STATIC_DEPS REQUIRED IMPORTED_TARGET
  freetype2
  harfbuzz
  libb2
  libbrotlidec
  libjpeg
  libpcre2-16
  libpng
  libtiff-4
  libwebp
  libwebpdemux
  libwebpmux
  openssl
  zlib
)

find_package(
  Qt6 ${AYANDICT_QT_VERSION} EXACT REQUIRED
  COMPONENTS Widgets Network MultimediaWidgets Svg
)

qt_add_executable(qt_link_probe main.cpp)
set_target_properties(qt_link_probe PROPERTIES WIN32_EXECUTABLE TRUE)
target_link_libraries(
  qt_link_probe
  PRIVATE
    Qt6::Widgets
    Qt6::Network
    Qt6::MultimediaWidgets
    PkgConfig::AYANDICT_STATIC_DEPS
)

# build-go.py passes the generated initializer/resource objects directly to the
# final Go external link. NO_DEFAULT keeps unrelated Qt plugins out.
qt_import_plugins(
  qt_link_probe
  NO_DEFAULT
  INCLUDE
    Qt6::QWindowsIntegrationPlugin
    Qt6::QModernWindowsStylePlugin
    Qt6::QWindowsMediaPlugin
    Qt6::QGifPlugin
    Qt6::QICOPlugin
    Qt6::QJpegPlugin
    Qt6::QSvgPlugin
    Qt6::QWebpPlugin
    Qt6::QTiffPlugin
)
