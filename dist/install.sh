#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<USAGE
Usage:
  $0 --tar <lightstar.tar.gz> [--target /] [--nodeps]

Options:
  --tar <file>      Path to lightstar tar/tar.gz package (optional when embedded)
  --target <dir>    Target root directory, default '/'
  --nodeps          Skip dependency installation
  -h, --help        Show this help
USAGE
}

ARCHIVE_FILE=""
TARGET_ROOT="/"
INSTALL_DEPS="yes"
SYS="linux"
TMP_DIR=""
INSTALLER="$0"
ARCHIVE_LINE=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --tar)
      ARCHIVE_FILE="${2:-}"
      shift 2
      ;;
    --target)
      TARGET_ROOT="${2:-}"
      shift 2
      ;;
    --nodeps|nodeps)
      INSTALL_DEPS="no"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$ARCHIVE_FILE" ]]; then
  ARCHIVE_LINE="$(grep -a -n "__ARCHIVE_BELOW__:$" "$INSTALLER" | cut -f1 -d: || true)"
  if [[ -z "$ARCHIVE_LINE" ]]; then
    echo "--tar is required when installer has no embedded archive" >&2
    usage
    exit 1
  fi
fi

if [[ -n "$ARCHIVE_FILE" && ! -f "$ARCHIVE_FILE" ]]; then
  echo "archive not found: $ARCHIVE_FILE" >&2
  exit 1
fi

if [[ ! -d "$TARGET_ROOT" ]]; then
  echo "target root not found: $TARGET_ROOT" >&2
  exit 1
fi

find_sys() {
  if command -v yum >/dev/null 2>&1; then
    SYS="redhat"
  elif command -v apt-get >/dev/null 2>&1; then
    SYS="debian"
  fi
}

requires() {
  echo "Installing dependencies ..."
  if [[ "$SYS" == "redhat" ]]; then
    yum update -y
    yum install -y libvirt-daemon libvirt qemu-kvm qemu-img unzip
  elif [[ "$SYS" == "debian" ]]; then
    apt-get update -y
    apt-get install -y libvirt-daemon-system libvirt-clients qemu-kvm qemu-utils unzip
  else
    echo "Unknown platform, skip dependency installation."
  fi
}

unpack_archive() {
  TMP_DIR="$(mktemp -d)"
  if [[ -n "$ARCHIVE_FILE" ]]; then
    echo "Unpacking archive: $ARCHIVE_FILE"
    tar -xf "$ARCHIVE_FILE" -C "$TMP_DIR" 2>/dev/null || tar -xzf "$ARCHIVE_FILE" -C "$TMP_DIR" || exit 1
    return
  fi

  local embedded_tar="$TMP_DIR/embedded.tar"
  echo "Unpacking embedded tar archive from installer ..."
  tail -n +"$((ARCHIVE_LINE + 1))" "$INSTALLER" > "$embedded_tar" || exit 1
  tar -xf "$embedded_tar" -C "$TMP_DIR" 2>/dev/null || tar -xzf "$embedded_tar" -C "$TMP_DIR" || exit 1
}

find_payload_root() {
  local roots
  roots="$(find "$TMP_DIR" -type d \( -name etc -o -name usr -o -name var \) -print | sed 's#/[a-z]*$##' | sort -u)"
  if [[ -z "$roots" ]]; then
    echo "Cannot find payload root containing etc/usr/var in zip" >&2
    exit 1
  fi
  echo "$roots" | head -n1
}

install_files() {
  local payload_root
  payload_root="$(find_payload_root)"

  echo "Installing files from: $payload_root"
  for d in etc usr var; do
    if [[ -d "$payload_root/$d" ]]; then
      mkdir -p "$TARGET_ROOT/$d"
      cp -af "$payload_root/$d/." "$TARGET_ROOT/$d/"
    fi
  done

  if [[ -d "$TARGET_ROOT/var/lightstar/script" ]]; then
    chmod +x "$TARGET_ROOT"/var/lightstar/script/*.sh || true
  fi

  if [[ -d "$payload_root" ]]; then
    find "$payload_root" -type f > "$TARGET_ROOT/usr/share/lightstar.db" 2>/dev/null || true
  fi
}

post_install() {
  if [[ "$TARGET_ROOT" != "/" ]]; then
    return
  fi

  if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload || true
    systemctl enable lightstar || true
  fi

  local etc_dir="/etc/lightstar"
  [ -e "$etc_dir/auth.json" ] || {
    cp -rf $etc_dir/auth.json.example $etc_dir/auth.json
    echo "!!! Please change default password for admin."
  }
  [ -e "$etc_dir/permission.json" ] || {
    cp -rf $etc_dir/permission.json.example $etc_dir/permission.json
  }
  [ -e "$etc_dir/zone.json" ] || {
    cat > $etc_dir/zone.json << EOF
{}
EOF
  }
}

cleanup() {
  [[ -n "$TMP_DIR" && -d "$TMP_DIR" ]] && rm -rf "$TMP_DIR"
}

trap cleanup EXIT

unpack_archive
find_sys
if [[ "$INSTALL_DEPS" == "yes" ]]; then
  requires
fi
install_files
post_install

echo "Install finished. target=${TARGET_ROOT}"

exit 0
