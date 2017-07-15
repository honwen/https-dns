FROM chenhw2/alpine:base
MAINTAINER CHENHW2 <https://github.com/chenhw2>

ARG VER=20170715
ARG URL=https://github.com/chenhw2/google-https-dns/releases/download/v$VER/google-https-dns_linux-amd64-$VER.tar.gz

RUN mkdir -p /usr/bin \
    && cd /usr/bin \
    && wget -qO- ${URL} | tar xz \
    && mv google-https-dns_* google-https-dns

USER nobody

ENV ARGS="-d 8.8.8.8"

EXPOSE 5300
EXPOSE 5300/udp

CMD google-https-dns -T -U ${ARGS} --logtostderr -V 3
