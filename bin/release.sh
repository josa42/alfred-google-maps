#!/bin/bash

version=$1

if [[ "$version" = "" ]]; then
  echo "usage: $0 <version>"
  exit 1
fi

rm -rf dist && mkdir dist

# glide install
go build main.go

zip -r dist/google-maps-$version.alfredworkflow . \
  -x vendor\* .git\* bin\* glide.yaml dist\* README.md glide.lock \*.go


hub release create -a dist/google-maps-$version.alfredworkflow $version