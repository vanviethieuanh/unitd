#!/usr/bin/env sh
set -eu

item="${1:-}"
if [ -z "$item" ]; then
    echo "ITEM path is required" >&2
    exit 1
fi

outdir="${2:-}"
if [ -z "$outdir" ]; then
    echo "Output directory is required" >&2
    exit 1
fi

PYTHON="${PYTHON:-.venv/bin/python3}"

mkdir -p tmp
system=$(printf '%s\n' "$item" | sed -E 's#.*/##; s/-gperf\.gperf(\.in)?$//')

# If the file is a Jinja2 template, expand it first.
if printf '%s\n' "$item" | grep -q '\.in$'; then
    input=$("$PYTHON" "$(dirname "$0")/expand_jinja2.py" "$item")
else
    input=$(cat "$item")
fi

printf '%s\n' "$input" |
rg '^[A-Za-z0-9_.]+,' |
    awk -F ',' '
{
  c1 = $1
  c2 = $2
  c3 = $3

  split(c1, a, ".")
  section=a[1]
  key=a[2]

  c4 = $4
  for (i = 5; i <= NF; i++)
  c4 = c4 "," $i

  gsub(/^[ \t]+|[ \t]+$/, "", c1)
  gsub(/^[ \t]+|[ \t]+$/, "", c2)
  gsub(/^[ \t]+|[ \t]+$/, "", c3)
  gsub(/^[ \t]+|[ \t]+$/, "", c4)

  print section "|" key "|" c2 "|" c3 "|" c4
}' |
    jq -c \
        --arg system "$system" \
        -R '
split("|") |
{
  system:  (if $system == "load-fragment" then "core" else $system end),
  section:  .[0],
  property: .[1],
  parser:   .[2],
  ltype:    .[3],
  offset:   .[4]
}
' >"$2/gperf_$system.jsonl"
