#!/bin/bash

set -e

OS="darwin linux windows"
ARCH="amd64 386"

echo "Ensuring code quality"
#go vet ./...
gofmt -w .

#ver=$(cd vendor/github.com/hashicorp/terraform && git describe --abbrev=0 --tags)
#echo "VERSION terraform '$ver'"

rm -Rf bin/
mkdir bin/

for GOOS in $OS; do
    for GOARCH in $ARCH; do
        arch="$GOOS-$GOARCH"
        binary="terraform-provider-sakuracloud"
        if [ "$GOOS" = "windows" ]; then
          binary="${binary}.exe"
        fi
        echo "Building $binary $arch"
        GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 govendor build -o $binary builtin/bins/provider-sakuracloud/main.go
        zip -r "bin/terraform-provider-sakuracloud_$arch" $binary
        rm -f $binary
    done
done
