FROM golang:1.20-alpine AS builder

ARG MIRROR=repo.huaweicloud.com
ARG TIMEZONE=Asia/Shanghai
ARG GOPROXY=https://goproxy.cn,https://proxy.golang.com.cn,direct
ARG GOPRIVATE=''
ARG PROJECT=''
ARG APP=''

RUN set -eux \
    && sed -i "s|dl-cdn.alpinelinux.org|$MIRROR|g" /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata make gcc musl-dev git upx \
    && cp -f /usr/share/zoneinfo/$TIMEZONE /etc/localtime \
    && echo "$TIMEZONE" > /etc/timezone

COPY . /srv/$PROJECT
WORKDIR /srv/$PROJECT

RUN set -eux \
    && ([ ! -z "$APP" ] || exit 1) \
    && make init \
    && make all \
    && make build/"$APP"

FROM alpine:latest

ARG MIRROR=repo.huaweicloud.com
ARG TIMEZONE=Asia/Shanghai
ARG PROJECT=''
ARG APP=''

RUN set -eux \
    && sed -i "s|dl-cdn.alpinelinux.org|$MIRROR|g" /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata \
    && cp -f /usr/share/zoneinfo/$TIMEZONE /etc/localtime \
    && echo "$TIMEZONE" > /etc/timezone

RUN set -eux \
    && ([ ! -z "$APP" ] || exit 1) \
    && mkdir -p /srv/configs

WORKDIR /srv

EXPOSE 8000
EXPOSE 9000

CMD ["/srv/app", "-config", "/srv/configs"]
