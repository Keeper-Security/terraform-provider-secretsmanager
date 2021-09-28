#!/bin/bash

#LATEST_TAG=`git describe --tags --abbrev=0`  #v0.1.2
LATEST_TAG=$1
VERSION=${LATEST_TAG#v}
VVERSION=v${VERSION}

echo LATEST_TAG=$LATEST_TAG

GOOS=darwin;  GOARCH=amd64 ;go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
GOOS=darwin;  GOARCH=arm64 ;go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
GOOS=freebsd; GOARCH=386   ;go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
GOOS=freebsd; GOARCH=amd64 ;go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
GOOS=freebsd; GOARCH=arm   ;go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
GOOS=linux;   GOARCH=386   ;go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
GOOS=linux;   GOARCH=amd64 ;go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
GOOS=linux;   GOARCH=arm   ;go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
GOOS=linux;   GOARCH=arm64 ;go build -o terraform-provider-keeper_${VVERSION}; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}
GOOS=windows; GOARCH=386   ;go build -o terraform-provider-keeper_${VVERSION}.exe; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}.exe
GOOS=windows; GOARCH=amd64 ;go build -o terraform-provider-keeper_${VVERSION}.exe; zip -m terraform-provider-keeper_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-keeper_${VVERSION}.exe
shasum -a 256 terraform-provider-keeper_${VERSION}_*.zip > terraform-provider-keeper_${VERSION}_SHA256SUMS
