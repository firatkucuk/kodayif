FROM ubuntu:14.04
RUN apt-get update && apt-get install -y git golang
RUN mkdir /opt/gopath && export GOPATH=/opt/gopath && go get github.com/streadway/amqp && git clone https://github.com/firatkucuk/kodayif.git /opt/kodayif
ADD entrypoint.sh /opt/entrypoint.sh
CMD /opt/entrypoint.sh
