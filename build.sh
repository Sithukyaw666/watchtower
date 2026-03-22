#!/bin/bash

BINFILE=watchtower
if [ -n "$MSYSTEM" ]; then
    BINFILE=watchtower.exe
fi
VERSION=$(git describe --tags)
echo "Building $VERSION for Linux (Static)..."

# Ensure we build a static Linux binary for the scratch container
CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-extldflags '-static' -X github.com/containrrr/watchtower/internal/meta.Version=$VERSION" -o $BINFILE .
