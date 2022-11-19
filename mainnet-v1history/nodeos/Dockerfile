FROM golang:latest AS health

RUN go install github.com/blockpane/fio-docker/health@latest

FROM ubuntu:18.04

RUN apt-get update; apt-get -y dist-upgrade; apt-get -y install aria2 jq pixz curl && \
    aria2c -x 5 "https://bin.fioprotocol.io/releases/fioprotocol-latest-ubuntu-18.04-amd64.deb" && \
    apt-get install -y ./fioprotocol-latest-ubuntu-18.04-amd64.deb && rm -f ./fioprotocol-latest-ubuntu-18.04-amd64.deb && \
    rm -f /etc/fio/nodeos/config.ini && rm -rf /var/lib/apt/lists/*
COPY fio-nodeos-run /usr/local/bin/fio-nodeos-run
COPY --from=health /go/bin/health /usr/local/bin/fio-health
RUN chmod 0755 /usr/local/bin/fio-nodeos-run

EXPOSE 8888/tcp
EXPOSE 3856/tcp
USER fio
WORKDIR /var/lib/fio

HEALTHCHECK --interval=10s --timeout=2s --start-period=10s CMD fio-health

CMD ["/usr/local/bin/fio-nodeos-run"]

