#!/bin/bash

# prepare variable
ENV=$CONTAINER_ENV
PRJ=`echo $CONTAINER_PROJ | awk -F"." '{print $NF}'`
IP=$CONTAINER_IP_ADDR
SCRIPT_STOP="{{ .StopCmd }}"

# stop
if [ -f /data/deploy/$PRJ/bin/$SCRIPT_STOP ] ;then
    /data/deploy/$PRJ/bin/$SCRIPT_STOP
fi