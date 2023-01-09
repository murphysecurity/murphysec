#!/usr/bin/env bash
set -e

if [ $(echo $CI_DEBUG_DELAY) ]; then
  echo Sleep for debug
  sleep 1h
fi

cd out/bin

proFile='murphysec-linux-amd64 murphysec-linux-arm64 murphysec-darwin-amd64 murphysec-windows-amd64.exe'

mkdir -p ../zip
zip pro.zip $proFile
mv pro.zip ../zip

cd ../zip

cos-uploader --local pro.zip --remote /client/$CI_BUILD_REF_NAME/pro.zip
cos-uploader --local pro.zip --remote /client/-/pro.zip
echo ================ SHA-256 ================
sha256sum *

