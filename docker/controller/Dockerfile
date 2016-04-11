FROM ubuntu:14.04
RUN apt-get update && apt-get install -y git golang
RUN mkdir /opt/gopath && export GOPATH=/opt/gopath && go get github.com/streadway/amqp && go get github.com/nu7hatch/gouuid && go get github.com/mediocregopher/radix.v2/redis && git clone https://github.com/firatkucuk/kodayif.git /opt/kodayif
ADD entrypoint.sh /opt/entrypoint.sh

EXPOSE 8080

CMD /opt/entrypoint.sh
