#!/bin/bash

cat <<EOF > /opt/kodayif/agent/config.json
{
  "controller"   : "${KODAYIF_CONTROLLER_PORT_8080_TCP_ADDR}:${KODAYIF_CONTROLLER_PORT_8080_TCP_PORT}",
  "mqConnString"  : "amqp://guest:guest@${RABBITMQ_PORT_5672_TCP_ADDR}:${RABBITMQ_PORT_5672_TCP_PORT}/"
}
EOF
cat /opt/kodayif/agent/config.json
export GOPATH=/opt/gopath
cd /opt/kodayif/agent/
go run *.go
