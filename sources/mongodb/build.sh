#!/bin/sh
# Compiles the dissector in this directory to WebAssembly.
#
#   ./build.sh            -> <directory name>.wasm, e.g. redis.wasm
#   ./build.sh myname     -> myname.wasm
#
# Uses TinyGo if it is installed, otherwise its official Docker image, so
# the only prerequisite is one of the two.
set -eu
cd "$(dirname "$0")"
NAME="${1:-redis}"
TINYGO_IMAGE="${TINYGO_IMAGE:-tinygo/tinygo:0.37.0}"
FLAGS="-target=wasm-unknown -opt=z -no-debug -panic=trap -scheduler=none"

if command -v tinygo >/dev/null 2>&1; then
  # shellcheck disable=SC2086
  tinygo build $FLAGS -o "$NAME.wasm" .
else
  # shellcheck disable=SC2086
  docker run --rm -v "$(pwd)":/src -w /src "$TINYGO_IMAGE" tinygo build $FLAGS -o "$NAME.wasm" .
fi
echo "built $NAME.wasm ($(wc -c < "$NAME.wasm" | tr -d ' ') bytes)"
echo "try it: tcpcat -sT -p <port> --scripts $(pwd) <target>"
