#!/usr/bin/env bash
set -ex

export CGO_ENABLED=0

ldflags="-s -w -X github.com/murphysecurity/murphysec/infra/buildinfo.version=$CI_COMMIT_TAG -buildid="

oss="windows linux darwin"
for i in  $oss; do
  export GOOS=$i
  outputName="murphysec-$i-amd64"
  if [ "$i" = 'windows' ]; then
      outputName="$outputName.exe"
  fi
  time go build -trimpath -ldflags "$ldflags" -o "out/bin/$outputName" ./cmd/murphy
done

