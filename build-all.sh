#!/bin/bash
# 
# Build different executable variables
set -e

outroot="build"

for os in windows linux darwin; do
  export GOOS="$os"

  ext=''
  if [[ "$os" == "windows" ]]; then
    ext='.exe'
  fi

  for arch in 386 amd64; do
    export GOARCH="$arch"

    outname="streamrelay-$os-$arch$ext"
    echo "$outname"
    go build -o "$outname"
  done
done