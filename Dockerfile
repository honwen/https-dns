FROM golang:alpine as builder
RUN apk add --update git
RUN go get github.com/chenhw2/google-https-dns

FROM chenhw2/alpine:base
MAINTAINER CHENHW2 <https://github.com/chenhw2>

# /usr/bin/google-https-dns
COPY --from=builder /go/bin /usr/bin

USER nobody

ENV ARGS="-d 8.8.8.8"

EXPOSE 5300
EXPOSE 5300/udp

CMD google-https-dns -T -U ${ARGS} --logtostderr -V 3
