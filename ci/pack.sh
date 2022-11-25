#!/usr/bin/env bash
set -e

if [ $(echo $CI_DEBUG_DELAY) ]; then
  echo Sleep for debug
  sleep 1h
fi

cd out/bin

saasFile='murphysec-saas-linux-amd64 murphysec-saas-darwin-amd64 murphysec-saas-windows-amd64.exe'
proFile='murphysec-linux-amd64 murphysec-darwin-amd64 murphysec-windows-amd64.exe'

mkdir -p ../zip
#zip saas.zip $saasFile
#mv saas.zip ../zip/
zip pro.zip $proFile
mv pro.zip ../zip

cd ../zip

cos-uploader --local pro.zip --remote /client/$CI_BUILD_REF_NAME/pro.zip
cos-uploader --local pro.zip --remote /client/-/pro.zip
#cos-uploader --local saas.zip --remote /client/$CI_BUILD_REF_NAME/saas.zip
#cos-uploader --local saas.zip --remote /client/-/saas.zip
echo ================ SHA-256 ================
sha256sum *

