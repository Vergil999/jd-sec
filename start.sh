#!/bin/bash
echo "商品ID:$1"
echo "接收登陆二维码邮箱:$2"
echo "抢购时间:$3"

./jdsec -itemId $1 -email $2 -secTime $3