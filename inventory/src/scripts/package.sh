#!/usr/bin/env bash
#
# Create a package tar of an install/source directory.
# Usage:
#   package.sh --install-src-dir /path/to/src --dist-dir /path/to/dist --version 1.2.3 [--app-name myapp]
#
# Environment variables can be used as fallback:
#   INSTALL_SRC_DIR, DIST_DIR, VERSION, APP_NAME
#
set -euo pipefail

# Defaults from env
INSTALL_SRC_DIR="${INSTALL_SRC_DIR:-}"
DIST_DIR="${DIST_DIR:-./dist}"
VERSION="${VERSION:-}"
APP_NAME="${APP_NAME:-}"

usage() {
  cat <<EOF
Usage:
  $0 --install-src-dir PATH --version VERSION [--dist-dir PATH] [--app-name NAME]

Parameters may also be supplied via env vars: INSTALL_SRC_DIR, DIST_DIR, VERSION, APP_NAME
Generated file: <app name>-v<sem version>-<os>-<arch>.tar (created in dist dir)
EOF
  exit 1
}

# Simple option parsing
while [[ $# -gt 0 ]]; do
  case "$1" in
    --install-src-dir) INSTALL_SRC_DIR="$2"; shift 2 ;;
    --dist-dir) DIST_DIR="$2"; shift 2 ;;
    --version) VERSION="$2"; shift 2 ;;
    --app-name) APP_NAME="$2"; shift 2 ;;
    -h|--help) usage ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      ;;
  esac
done

# Validate inputs
if [[ -z "${INSTALL_SRC_DIR:-}" ]]; then
  echo "ERROR: install source dir not specified" >&2
  usage
fi
if [[ -z "${VERSION:-}" ]]; then
  echo "ERROR: version not specified" >&2
  usage
fi
if [[ ! -d "$INSTALL_SRC_DIR" ]]; then
  echo "ERROR: install source dir does not exist or is not a directory: $INSTALL_SRC_DIR" >&2
  exit 2
fi

# Determine app name
if [[ -z "${APP_NAME}" ]]; then
  # default to basename of install src dir
  APP_NAME="$(basename -s / "$INSTALL_SRC_DIR")"
fi

# Normalize version: remove leading 'v' if present, then prefix later
VERSION_STRIPPED="${VERSION#v}"
# basic semver-ish validation (very permissive)
if ! [[ "$VERSION_STRIPPED" =~ ^[0-9]+(\.[0-9]+){0,2}([.-].*)?$ ]]; then
  echo "Warning: version '$VERSION' does not look like a semantic version. Proceeding anyway." >&2
fi

# Detect OS and arch
_uname_s="$(uname -s 2>/dev/null || echo unknown)"
_uname_m="$(uname -m 2>/dev/null || echo unknown)"

# normalize os
# POSIX-compatible lowercase conversion to support older bash (e.g. macOS bash 3.2)
uname_s_lc="$(printf '%s' "$_uname_s" | tr '[:upper:]' '[:lower:]')"
case "${uname_s_lc}" in
  linux*) os="linux" ;;
  darwin*) os="darwin" ;;
  mingw*|msys*|cygwin*) os="windows" ;;
  *) os="${uname_s_lc}" ;;
esac

# normalize arch
case "${_uname_m}" in
  x86_64|amd64) arch="amd64" ;;
  aarch64|arm64) arch="arm64" ;;
  armv7l|armv7) arch="armv7" ;;
  i386|i686) arch="386" ;;
  *) arch="${_uname_m}" ;;
esac

# Prepare output paths
dist_dir="$DIST_DIR"
mkdir -p "$dist_dir"
chmod 0755 "$dist_dir" || true

pkg_name="${APP_NAME}-v${VERSION_STRIPPED}-${os}-${arch}.tar"
dest_path="$(cd "$dist_dir" && pwd)/${pkg_name}"

# create staging dir and copy contents into a top-level folder named "<app>-v<version>"
staging="$(mktemp -d)"
trap 'rm -rf -- "$staging"' EXIT

topdir="${APP_NAME}-v${VERSION_STRIPPED}"
mkdir -p "$staging/$topdir"

# copy contents (preserve attributes)
if command -v rsync >/dev/null 2>&1; then
  rsync -a --delete "$INSTALL_SRC_DIR"/ "$staging/$topdir"/
else
  cp -a "$INSTALL_SRC_DIR"/. "$staging/$topdir"/
fi

# create tar (no compression, as requested .tar)
(
  cd "$staging"
  tar -cf "$dest_path" "$topdir"
)

# set permissions on the tar (readable)
chmod 0644 "$dest_path" || true

echo "$dest_path"