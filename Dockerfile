# syntax=docker/dockerfile:1

FROM curlimages/curl:7.81.0 AS downloader
WORKDIR /tmp
RUN VERSION=v2.3.1 && \
  curl -sSL -o argocd https://github.com/argoproj/argo-cd/releases/download/$VERSION/argocd-linux-amd64 && \
  chmod +x argocd

FROM golang:1.16-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY src ./
RUN GOOS=linux go build -o plugin

FROM alpine:3.15
WORKDIR /app
COPY --from=downloader /tmp/argocd /usr/local/bin/argocd
COPY --from=build /tmp/plugin ./plugin

RUN useradd -u 8737 argo && \
  chowm argo /app/bin/plugin
USER argo

ENTRYPOINT [ "/app/bin/plugin" ]

