#!/bin/bash

version=$1

if [[ "${version}" = "" ]]; then
  echo "usage: ${0} <version>"
  exit 1
fi

rm -rf dist && mkdir dist

# glide install
dep ensure || exit 1
go build main.go

defaults write "$(pwd)/info.plist" version "${version}"
plutil -convert xml1  "$(pwd)/info.plist"
git add info.plist

cat README.md | sed 's/download\/\([^\/]*\)\/google-maps-.*\.alfredworkflow/download\/'$version'\/google-maps-'$version'.alfredworkflow/' > README.md.new
mv -f README.md.new README.md
git add README.md

git cm "🎉  Release ${version}"
git push

zip -r "dist/google-maps-${version}.alfredworkflow" . \
  -x vendor\* .git\* bin\* Gopkg.\* dist\* README.md glide.lock \*.go \*.DS_Store docs/\*

# git tag "${version}" && git push --tags

hub release create \
  -m "🎉  Release ${version}" \
  -a "dist/google-maps-$version.alfredworkflow" \
  "${version}"
