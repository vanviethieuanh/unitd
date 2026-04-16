#!/usr/bin/env sh
set -eu

name="${1:-}"
if [ -z "$name" ]; then
  echo "NAME is required (example: task download-man NAME=systemd.timer)" >&2
  exit 1
fi

url="https://raw.githubusercontent.com/systemd/systemd/main/man/${name}.xml"

out_dir="./tmp/man"
etag_dir="./tmp/man/etag"
out_file="${out_dir}/${name}.xml"
etag_file="${etag_dir}/${name}.xml.etag"
tmp_file="${out_file}.tmp"

mkdir -p "$out_dir" "$etag_dir"

headers=$(mktemp)
cleanup() {
  rm -f "$headers"
}
trap cleanup EXIT

if [ -f "$etag_file" ]; then
  http_code=$(curl -sSL -D "$headers" -H "If-None-Match: $(cat "$etag_file")" -o "$tmp_file" -w '%{http_code}' "$url")
else
  http_code=$(curl -sSL -D "$headers" -o "$tmp_file" -w '%{http_code}' "$url")
fi

case "$http_code" in
  200)
    mv "$tmp_file" "$out_file"
    etag=$(awk 'BEGIN{IGNORECASE=1} /^etag:/ {sub(/\r$/, "", $0); sub(/^[Ee][Tt][Aa][Gg]:[[:space:]]*/, "", $0); print; exit}' "$headers")
    if [ -n "$etag" ]; then
      printf '%s\n' "$etag" > "$etag_file"
    fi
    echo "Downloaded: $out_file"
    ;;
  304)
    rm -f "$tmp_file"
    echo "Not modified: $out_file"
    ;;
  *)
    rm -f "$tmp_file"
    echo "Failed to download $url (HTTP $http_code)" >&2
    exit 1
    ;;
esac
