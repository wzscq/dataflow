#!/bin/sh
echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

if [ ! -e package/web ]; then
  mkdir package/web
fi

echo build the code ...
cd ../flowconfig
npm install
sed -i  's/host=\"*.*\"/host=\"\"/' ./public/index.html
npm run build
cd ../build

echo remove last package if exist
if [ -e package/web/flowconfig ]; then
  rm -rf package/web/flowconfig
fi

mv ../flowconfig/build ./package/web/flowconfig

echo flowconfig package build over.
