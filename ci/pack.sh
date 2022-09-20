#!/usr/bin/env bash
set -e
cd out
if [ $(echo $CI_DEBUG_DELAY) ]; then
  echo Sleep for debug
  sleep 1h
fi

saasFile='murphysec-saas-linux-amd64 murphysec-saas-darwin-amd64 murphysec-saas-windows-amd64.exe'
proFile='murphysec-linux-amd64 murphysec-darwin-amd64 murphysec-windows-amd64.exe'

zip saas.zip $saasFile
zip pro.zip $proFile
echo -n $saasFile | xargs -I % -d ' ' cos-uploader --local % --remote /client/$CI_BUILD_REF_NAME/%
echo -n $saasFile | xargs -I % -d ' ' cos-uploader --local % --remote /client/-/%
echo -n $proFile | xargs -I % -d ' ' cos-uploader --local % --remote /client/$CI_BUILD_REF_NAME/%
echo -n $proFile | xargs -I % -d ' ' cos-uploader --local % --remote /client/-/%
cos-uploader --local pro.zip --remote /client/$CI_BUILD_REF_NAME/pro.zip
cos-uploader --local pro.zip --remote /client/-/pro.zip
cos-uploader --local saas.zip --remote /client/$CI_BUILD_REF_NAME/saas.zip
cos-uploader --local saas.zip --remote /client/-/saas.zip
echo ================ SHA-256 ================
sha256sum *
