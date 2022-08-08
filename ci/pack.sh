#!/usr/bin/env bash
set -e
cd out
if [ $(echo $CI_DEBUG_DELAY) ]; then
  echo Sleep for debug
  sleep 1h
fi

saasFile='murphysec-saas-linux-amd64 murphysec-saas-darwin-amd64 murphysec-saas-windows-amd64.exe'
proFile='murphysec-linux-amd64 murphysec-darwin-amd64 murphysec-windows-amd64.exe'

if [ $(echo $CI_BUILD_REF_NAME | grep -o '^v') ]; then
  echo Packaging for IDEA Plugin
  f_prefix=$(echo $CI_BUILD_REF_NAME | sed -r 's/(^v|\.)//g')
  echo Name prefix: $f_prefix
  cp murphysec-saas-linux-amd64 $f_prefix-murphysec-saas-linux-amd64
  cp murphysec-saas-darwin-amd64 $f_prefix-murphysec-saas-darwin-amd64
  cp murphysec-saas-windows-amd64.exe $f_prefix-murphysec-saas-windows-amd64.exe
  ideaFile="$f_prefix-murphysec-saas-linux-amd64 $f_prefix-murphysec-saas-darwin-amd64 $f_prefix-murphysec-saas-windows-amd64.exe"
  zip idea.zip $ideaFile
  echo Uploading...
  cos-uploader --local idea.zip --remote /client/$CI_BUILD_REF_NAME/idea.zip
  cos-uploader --local idea.zip --remote /client/-/idea.zip
fi

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
