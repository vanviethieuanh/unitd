#!/usr/bin/env sh
set -eu

cat $1/src/core/load-fragment.c |
    grep '^[[:space:]]*{[[:space:]]*config_parse' |
    sed -E 's/[{}]//g' |
    awk -F ',' '{
    gsub(/^[[:space:]]+|[[:space:]]+$/, "", $1)
    gsub(/^[[:space:]]+|[[:space:]]+$/, "", $2)
    gsub(/"/, "", $2)
    print $1 "," $2
}' |
    jq -c \
        -R '
split(",") |
{
  parser: .[0],
  type: .[1]
}
' >$2
