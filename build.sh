#!/bin/bash

# List of platforms to build for
platforms=("darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64" "windows/amd64" "windows/arm64")

for platform in "${platforms[@]}"
do
  IFS="/" read -r GOOS GOARCH <<< "$platform"
  output_name="nativeblocks-$GOOS-$GOARCH"
  
  # Add .exe extension for Windows binaries
  if [ "$GOOS" = "windows" ]; then
    output_name+=".exe"
  fi
  
  echo "Building $output_name"
  GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name
done
