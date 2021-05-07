#!/bin/bash

# prepare variable
ENV=$CONTAINER_ENV
PRJ=`echo $CONTAINER_PROJ | awk -F"." '{print $NF}'`
IP=$CONTAINER_IP_ADDR
PORT="{{ .HTTPPort }}"
RPCPORT="{{ .RPCPort }}"
SCRIPT_START="{{ .StartCmd }}"

# prepare env
mkdir -p /persist/logs
if [ ! -h /data/deploy/$PRJ/logs ];then
  ln -s /persist/logs /data/deploy/$PRJ/logs
fi

{{ .PreCmd }}

# start
cd /data/deploy/$PRJ
nohup $SCRIPT_START >> logs/nohup.out 2>&1 &
echo $! > /persist/logs/app.pid