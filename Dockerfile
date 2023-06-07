# API build stage
FROM golang:1.20.4-alpine3.18 as go-builder
ARG GOPROXY=goproxy.cn

ENV GOPROXY=https://${GOPROXY},direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache make bash git tzdata

WORKDIR /ecodepost

COPY go.mod go.sum ./
RUN go mod download -x
COPY scripts scripts
COPY bff bff
COPY user-svc user-svc
COPY resource-svc resource-svc
COPY job job
COPY proto proto
COPY pb pb
COPY sdk sdk
COPY main.go main.go
COPY config config
COPY Makefile Makefile
RUN ls -rlt ./bff/pkg/server/ui/dist && make build.api
RUN ls -rlt /ecodepost


# Fianl running stage
FROM alpine:3.14.3
LABEL maintainer="ecodepost@ecodeclub.member"

WORKDIR /ecodepost

COPY --from=go-builder /ecodepost/../bin/ecodepost ./bin/
COPY --from=go-builder /ecodepost/config ./config

EXPOSE 9002

RUN apk add --no-cache tzdata

CMD ["sh", "-c", "./bin/ecodepost"]
