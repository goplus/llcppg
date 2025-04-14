#!/bin/bash
set -e

sudo apt install liblua5.4-dev libsqlite3-dev libgmp-dev libgpg-error-dev zlib1g-dev libisl-dev libz3-dev -y
llcppgtest -demo ./_llcppgtest/cjson -conf conf/linux
llcppgtest -demo ./_llcppgtest/gmp -conf conf/linux
llcppgtest -demo ./_llcppgtest/gpgerror -conf conf/linux
llcppgtest -demo ./_llcppgtest/isl
llcppgtest -demo ./_llcppgtest/lua -conf conf/linux
llcppgtest -demo ./_llcppgtest/sqlite -conf conf/linux
llcppgtest -demo ./_llcppgtest/z3 -conf conf/linux
llcppgtest -demo ./_llcppgtest/zlib -conf conf/linux
