#!/bin/bash

cat <<EOF > /opt/kodayif/controller/config.json
{
  "listenAddress" : ":8080",
  "mqConnString"  : "amqp://guest:guest@${RABBITMQ_PORT_5672_TCP_ADDR}:${RABBITMQ_PORT_5672_TCP_PORT}/",
  "redisAddress"  : "${REDIS_PORT_6379_TCP_ADDR}:${REDIS_PORT_6379_TCP_PORT}"
}
EOF
cat /opt/kodayif/controller/config.json
export GOPATH=/opt/gopath
cd /opt/kodayif/controller/
go run *.go
