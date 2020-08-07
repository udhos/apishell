#!/bin/bash

export CGO_ENABLED=0
export GO111MODULE=on

build() {
	local pkg="$1"

	gofmt -s -w "$pkg"
	go fix "$pkg"
	go vet "$pkg"

	#hash gosimple >/dev/null && gosimple "$pkg"
	hash golint >/dev/null && golint "$pkg"
	#hash staticcheck >/dev/null && staticcheck "$pkg"

	go test -failfast "$pkg"
	go install -v "$pkg"
}

build ./apid
build ./apictl
