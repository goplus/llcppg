#!/bin/bash
set -e

# for test
go install ./cmd/llcppcfg
go install ./cmd/llcppgtest

# main process required
llgo install ./_xtool/llcppsymg
llgo install ./_xtool/llcppsigfetch
go install ./cmd/gogensig
go install .
