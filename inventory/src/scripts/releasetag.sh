#!/usr/bin/env bash
# generate_release_tag VERSION
# If VERSION is omitted, uses $VERSION from environment.
generate_release_tag() {
  local input VERSION_STRIPPED _uname_s _uname_m uname_s_lc os arch

  input="${1:-${VERSION:-}}"
  if [ -z "$input" ]; then
    echo "generate_release_tag: VERSION not provided" >&2
    return 2
  fi

  # Normalize version: remove leading 'v' if present, then prefix later
  VERSION_STRIPPED="${input#v}"

  # basic semver-ish validation (very permissive)
  if ! [[ "$VERSION_STRIPPED" =~ ^[0-9]+(\.[0-9]+){0,2}([.-].*)?$ ]]; then
    echo "Warning: version '$input' does not look like a semantic version. Proceeding anyway." >&2
  fi

  # Detect OS and arch
  _uname_s="$(uname -s 2>/dev/null || echo unknown)"
  _uname_m="$(uname -m 2>/dev/null || echo unknown)"

  # normalize os (POSIX-compatible lowercase conversion)
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

  # Output tag: v${VERSION_STRIPPED}-${os}-${arch}
  printf 'v%s-%s-%s' "$VERSION_STRIPPED" "$os" "$arch"
}

version_stripped() {
  input="${1:-${VERSION:-}}"

  # Normalize version: remove leading 'v' if present, then prefix later
  VERSION_STRIPPED="${input#v}"

  # basic semver-ish validation (very permissive)
  if ! [[ "$VERSION_STRIPPED" =~ ^[0-9]+(\.[0-9]+){0,2}([.-].*)?$ ]]; then
    echo "Warning: version '$input' does not look like a semantic version. Proceeding anyway." >&2
  fi
  echo "$VERSION_STRIPPED"
}