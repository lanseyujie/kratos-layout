FROM golang:1.20-alpine AS builder

ARG MIRROR=repo.huaweicloud.com
ARG TIMEZONE=Asia/Shanghai
ARG GOPROXY=https://goproxy.cn,https://proxy.golang.com.cn,direct
ARG GOPRIVATE=''
ARG PRIVATEKEY=''
ARG PROJECT=''
ARG APP=''

RUN set -eux \
    && sed -i "s|dl-cdn.alpinelinux.org|$MIRROR|g" /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata make gcc musl-dev upx openssh-client git \
    && cp -f /usr/share/zoneinfo/$TIMEZONE /etc/localtime \
    && echo "$TIMEZONE" > /etc/timezone

RUN set -eux \
    && [ ! -z "$GOPRIVATE" ] \
    && [ ! -z "$PRIVATEKEY" ] \
    && cd ~ \
    && mkdir -p ~/.ssh/ \
    && echo -e "Host *\n\tAddKeysToAgent yes\n\tStrictHostKeyChecking no\n" >~/.ssh/config \
    && echo -n "$PRIVATEKEY" | base64 -d > ~/.ssh/id_rsa \
    && chmod 600 ~/.ssh/id_rsa \
    && ssh -T -v "git@$GOPRIVATE" \
    && git config --global url."ssh://git@$GOPRIVATE/".insteadOf "https://$GOPRIVATE/" || true

COPY . /srv/$PROJECT
WORKDIR /srv/$PROJECT

RUN set -eux \
    && [ ! -z "$APP" ] \
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

COPY --from=builder /srv/$PROJECT/bin/$APP /srv/app

WORKDIR /srv

EXPOSE 8000
EXPOSE 9000

CMD ["/srv/app", "-config", "/srv/configs"]
