#!/bin/bash
set -e

brew install lua zlib isl libgpg-error raylib z3 sqlite3 gmp libxml2 libxslt

export PKG_CONFIG_PATH="/opt/homebrew/opt/zlib/lib/pkgconfig"
export PKG_CONFIG_PATH="/opt/homebrew/opt/sqlite/lib/pkgconfig:$PKG_CONFIG_PATH"
export PKG_CONFIG_PATH="/opt/homebrew/opt/libxml2/lib/pkgconfig:$PKG_CONFIG_PATH"
export PKG_CONFIG_PATH="/opt/homebrew/opt/libxslt/lib/pkgconfig:$PKG_CONFIG_PATH"
pkg-config --cflags --libs sqlite3
pkg-config --cflags --libs libxslt

llcppgtest -demos ./_llcppgtest