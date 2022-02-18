#!/bin/bash
git pull origin master
echo '-------------拉取最新代码over--------------'
echo '-------------打包镜像--------------'
docker build -t liyongqi666/jd-sec:1.0 .
echo '-------------打包完毕--------------'


