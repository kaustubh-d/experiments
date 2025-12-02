#!/usr/bin/env bash
set -euo pipefail

# Defaults (can be overridden via environment or CLI flags)
INSTALL_DIR="${INSTALL_DIR:-/usr/local}"
BUILD_DIR="${BUILD_DIR:-./build}"
CONFIG_DIR="${CONFIG_DIR:-./config}"

usage() {
  cat <<EOF
Usage: $0 [--build-dir DIR] [--config-dir DIR] [--install-dir DIR]

Environment variables also supported: BUILD_DIR, CONFIG_DIR, INSTALL_DIR
This script copies:
 - binaries from BUILD_DIR/bin -> INSTALL_DIR/bin (made executable)
 - config files from CONFIG_DIR (or BUILD_DIR/config) -> INSTALL_DIR/config (0644)
EOF
  exit 1
}

# Simple CLI parsing
while [ "$#" -gt 0 ]; do
  case "$1" in
    --build-dir=*) BUILD_DIR="${1#*=}"; shift ;;
    --build-dir) BUILD_DIR="$2"; shift 2 ;;
    -b) BUILD_DIR="$2"; shift 2 ;;
    --config-dir=*) CONFIG_DIR="${1#*=}"; shift ;;
    --config-dir) CONFIG_DIR="$2"; shift 2 ;;
    -c) CONFIG_DIR="$2"; shift 2 ;;
    --install-dir=*) INSTALL_DIR="${1#*=}"; shift ;;
    --install-dir) INSTALL_DIR="$2"; shift 2 ;;
    -i) INSTALL_DIR="$2"; shift 2 ;;
    -h|--help) usage ;;
    *) echo "Unknown argument: $1"; usage ;;
  esac
done

BIN_SRC="$BUILD_DIR/bin"
BIN_DST="$INSTALL_DIR/bin"
CONF_SRC="$CONFIG_DIR"
CONF_DST="$INSTALL_DIR/config"

mkdir -p "$BIN_DST" "$CONF_DST"

echo "Installing binaries from '$BIN_SRC' -> '$BIN_DST'..."
if compgen -G "$BIN_SRC/*" > /dev/null 2>&1; then
  for src in "$BIN_SRC"/*; do
    [ -f "$src" ] || continue
    cp -f -- "$src" "$BIN_DST"/
    chmod 0755 "$BIN_DST/$(basename "$src")"
    echo " installed $BIN_DST/$(basename "$src")"
  done
else
  echo "No binaries found in '$BIN_SRC' (skipping)"
fi

echo "Installing config files from '$CONF_SRC' -> '$CONF_DST'..."
if compgen -G "$CONF_SRC/*" > /dev/null 2>&1; then
  for src in "$CONF_SRC"/*; do
    [ -f "$src" ] || continue
    cp -f -- "$src" "$CONF_DST"/
    chmod 0644 "$CONF_DST/$(basename "$src")"
    echo " installed $CONF_DST/$(basename "$src")"
  done
else
  echo "No config files found in '$CONF_SRC' (skipping)"
fi