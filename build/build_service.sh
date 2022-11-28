#!/bin/sh
echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

echo build the code ...
cd ../service
#添加参数CGO_ENABLED=0，关闭CGO,是为了是编译后的程序可以在alpine中运行
CGO_ENABLED=0 go build -o dataflow
cd ../build

echo remove last package if exist
if [ ! -e package/service ]; then
  mkdir package/service
fi

if [ -e package/service/dataflow ]; then
  rm -rf package/service/dataflow
fi

mv ../service/dataflow ./package/service/dataflow

echo service dataflow build over.
