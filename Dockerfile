# UI build stage
#FROM node:16-alpine3.14 as js-builder
#
#ENV NODE_OPTIONS=--max_old_space_size=8000
#WORKDIR /ecodepost
#COPY ecodepost-fe .
#WORKDIR /ecodepost/ecodepost-fe
#RUN yarn install --frozen-lockfile
#RUN npm run build
#RUN cp -rf ./dist /ecodepost/dist
#RUN cd ../ && rm -rf ecodepost-fe


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


# Fianl running stage
FROM alpine:3.14.3
LABEL maintainer="ecodepost@ecodeclub.member"

WORKDIR /ecodepost

COPY --from=go-builder /ecodepost/bin/clickvisual ./bin/
COPY --from=go-builder /ecodepost/config ./config

EXPOSE 9002

RUN apk add --no-cache tzdata

CMD ["sh", "-c", "./bin/ecodepost"]