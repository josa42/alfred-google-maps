#!/bin/bash

version=$1

if [[ "${version}" = "" ]]; then
  echo "usage: ${0} <version>"
  exit 1
fi

rm -rf dist && mkdir dist

# glide install
go build main.go

defaults write "$(pwd)/info.plist" version "${version}"
plutil -convert xml1  "$(pwd)/info.plist"

git add info.plist
git cm "🎉  Release ${version}"
git push

zip -r "dist/google-maps-${version}.alfredworkflow" . \
  -x vendor\* .git\* bin\* glide.yaml dist\* README.md glide.lock \*.go

# git tag "${version}" && git push --tags

hub release create \
  -m "🎉  Release ${version}" \
  -a "dist/google-maps-$version.alfredworkflow" \
  "${version}"