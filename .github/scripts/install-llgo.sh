#!/usr/bin/env bash
# Install LLGo (plus the LLVM and Go it needs) so that `llgo test` works.
#
# Usage: .github/scripts/install-llgo.sh
#
# Environment overrides (defaults are what llcppg CI uses):
#   LLGO_REF       llgo git ref (branch, tag or commit)   default: main
#   LLVM_VERSION   LLVM major version                     default: 22
#   LLGO_SRC       where llgo sources are checked out     default: $HOME/.llgo-src
#   LLGO_BIN_DIR   where the `llgo` binary is installed   default: $HOME/.llgo/bin
#   GO_VERSION     Go version to use                      default: `go` line of go.mod
#   GO_ROOT        where a downloaded Go is installed     default: $HOME/.llgo/go
#
# Supports Ubuntu/Debian (apt) and macOS (Homebrew). Safe to re-run.
# After it finishes: export PATH="$LLGO_BIN_DIR:$PATH" (also printed at the end).
set -euo pipefail

LLGO_REF="${LLGO_REF:-main}"
LLVM_VERSION="${LLVM_VERSION:-22}"
LLGO_SRC="${LLGO_SRC:-$HOME/.llgo-src}"
LLGO_BIN_DIR="${LLGO_BIN_DIR:-$HOME/.llgo/bin}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

log() { echo "==> $*"; }
SUDO=""
if [ "$(id -u)" -ne 0 ]; then SUDO="sudo"; fi

install_llvm_linux() {
  source "$REPO_ROOT/.github/scripts/preserve-extra-ca.sh"
  preserve_extra_ca "$SUDO"

  if [ -x "/usr/lib/llvm-$LLVM_VERSION/bin/llvm-config" ]; then
    log "LLVM $LLVM_VERSION already installed"
  else
    local codename
    codename="$(. /etc/os-release && echo "${VERSION_CODENAME:-$(lsb_release -cs)}")"
    $SUDO apt-get update
    $SUDO apt-get install -y wget gnupg ca-certificates lsb-release
    wget -qO- https://apt.llvm.org/llvm-snapshot.gpg.key |
      $SUDO tee /etc/apt/trusted.gpg.d/apt.llvm.org.asc >/dev/null
    echo "deb https://apt.llvm.org/$codename/ llvm-toolchain-$codename-$LLVM_VERSION main" |
      $SUDO tee /etc/apt/sources.list.d/llvm-$LLVM_VERSION.list >/dev/null
    $SUDO apt-get update
    $SUDO apt-get install -y "llvm-$LLVM_VERSION-dev" "clang-$LLVM_VERSION" \
      "libclang-$LLVM_VERSION-dev" "lld-$LLVM_VERSION" \
      "libunwind-$LLVM_VERSION-dev" "libc++-$LLVM_VERSION-dev"
  fi
  $SUDO apt-get install -y git curl cmake pkg-config libgc-dev libssl-dev \
    zlib1g-dev libffi-dev libuv1-dev
  LLVM_BIN="/usr/lib/llvm-$LLVM_VERSION/bin"
}

install_llvm_macos() {
  brew install "llvm@$LLVM_VERSION" "lld@$LLVM_VERSION" bdw-gc openssl libffi libuv pkg-config
  brew link --force --overwrite "llvm@$LLVM_VERSION" "lld@$LLVM_VERSION" || true
  LLVM_BIN="$(brew --prefix "llvm@$LLVM_VERSION")/bin"
  export PATH="$(brew --prefix "lld@$LLVM_VERSION")/bin:$PATH"
}

case "$(uname -s)" in
  Linux)  install_llvm_linux ;;
  Darwin) install_llvm_macos ;;
  *) echo "unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac
export PATH="$LLVM_BIN:$PATH"

# Go: llgo shells out to the `go` on PATH, and that launcher itself must be the
# required version (an older launcher rejects newer GOEXPERIMENT values, e.g.
# "unknown GOEXPERIMENT dwarf5"). So the real Go, not just a toolchain-switching
# stub, has to be >= GO_VERSION (taken from go.mod).
GO_VERSION="${GO_VERSION:-$(awk '/^go /{print $2; exit}' "$REPO_ROOT/go.mod")}"
GO_ROOT="${GO_ROOT:-$HOME/.llgo/go}" # used only when downloading
go_ok() {
  command -v go >/dev/null &&
    [ "$(GOTOOLCHAIN=local go env GOVERSION 2>/dev/null)" = "go$GO_VERSION" ]
}
if [ -x "$GO_ROOT/bin/go" ]; then export PATH="$GO_ROOT/bin:$PATH"; fi
if ! go_ok && command -v go >/dev/null; then
  # Any Go >= 1.21 can fetch the exact toolchain through the Go module proxy;
  # put that toolchain's own bin directory first on PATH.
  log "Fetching Go $GO_VERSION toolchain via the Go module proxy"
  if toolchain_root="$(GOTOOLCHAIN="go$GO_VERSION" go env GOROOT)" && [ -x "$toolchain_root/bin/go" ]; then
    GO_ROOT="$toolchain_root"
    export PATH="$GO_ROOT/bin:$PATH"
  fi
fi
if ! go_ok; then
  case "$(uname -s)-$(uname -m)" in
    Linux-x86_64) platform=linux-amd64 ;;
    Linux-aarch64|Linux-arm64) platform=linux-arm64 ;;
    Darwin-arm64) platform=darwin-arm64 ;;
    Darwin-x86_64) platform=darwin-amd64 ;;
    *) echo "unsupported platform" >&2; exit 1 ;;
  esac
  log "Downloading Go $GO_VERSION into $GO_ROOT"
  rm -rf "$GO_ROOT" "$GO_ROOT.tmp"
  mkdir -p "$GO_ROOT.tmp"
  curl -fsSL "https://go.dev/dl/go$GO_VERSION.$platform.tar.gz" |
    tar -C "$GO_ROOT.tmp" -xzf -
  mv "$GO_ROOT.tmp/go" "$GO_ROOT"
  rmdir "$GO_ROOT.tmp"
  export PATH="$GO_ROOT/bin:$PATH"
fi
export GOTOOLCHAIN=local
log "Using $(go version)"

log "Fetching llgo ($LLGO_REF) into $LLGO_SRC"
if [ ! -d "$LLGO_SRC/.git" ]; then
  git clone https://github.com/goplus/llgo "$LLGO_SRC"
fi
git -C "$LLGO_SRC" fetch --tags origin
git -C "$LLGO_SRC" checkout --force "$LLGO_REF" 2>/dev/null ||
  git -C "$LLGO_SRC" checkout --force "origin/$LLGO_REF"
if git -C "$LLGO_SRC" symbolic-ref -q HEAD >/dev/null; then
  git -C "$LLGO_SRC" reset --hard "origin/$LLGO_REF"
fi

log "Building llgo"
mkdir -p "$LLGO_BIN_DIR"
(
  cd "$LLGO_SRC"
  export CGO_CPPFLAGS="$(llvm-config --cflags)"
  export CGO_CXXFLAGS="$(llvm-config --cxxflags)"
  export CGO_LDFLAGS="$(llvm-config --ldflags --libs --system-libs)"
  go build -tags=dev,byollvm -o "$LLGO_BIN_DIR/llgo" ./cmd/llgo
)

export PATH="$LLGO_BIN_DIR:$PATH"
export LLGO_ROOT="$LLGO_SRC"
log "llgo version"
llgo version

if [ -n "${GITHUB_PATH:-}" ]; then
  echo "$LLGO_BIN_DIR" >> "$GITHUB_PATH"
  echo "$GO_ROOT/bin" >> "$GITHUB_PATH"
  echo "$LLVM_BIN" >> "$GITHUB_PATH"
fi
if [ -n "${GITHUB_ENV:-}" ]; then
  echo "LLGO_ROOT=$LLGO_SRC" >> "$GITHUB_ENV"
  echo "GOTOOLCHAIN=local" >> "$GITHUB_ENV"
fi

cat <<MSG

LLGo installed. In your shell run:
  export PATH="$LLGO_BIN_DIR:$LLVM_BIN:$GO_ROOT/bin:\$PATH"
  export GOTOOLCHAIN=local
  export LLGO_ROOT="$LLGO_SRC"
Then verify with: llgo version && (cd "$REPO_ROOT" && llgo test ./tool/pputil/...)
MSG
