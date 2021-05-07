#!/bin/bash

# prepare variable
ENV=$CONTAINER_ENV
PRJ=`echo $CONTAINER_PROJ | awk -F"." '{print $NF}'`
IP=$CONTAINER_IP_ADDR
PORT="{{ .HTTPPort }}"
RPCPORT="{{ .RPCPort }}"
SCRIPT_START="{{ .StartCmd }}"

# prepare env
export _JAVA_OPTIONS="-Djava.net.preferIPv4Stack=true -Dfile.encoding=UTF-8 -Dxueqiu.env=$ENV -Dxueqiu.service=$PRJ -Dxueqiu.ip=$IP"
mkdir -p /persist/logs
if [ ! -h /data/deploy/$PRJ/logs ];then
  ln -s /persist/logs /data/deploy/$PRJ/logs
fi

{{ .PreCmd }}

# start
cd /data/deploy/$PRJ
nohup bin/$SCRIPT_START --env=$ENV --ip=$IP --port=$PORT >> logs/nohup.out 2>&1 &
echo $! > /persist/logs/app.pid