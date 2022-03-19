# syntax=docker/dockerfile:1

FROM curlimages/curl:7.81.0 AS downloader
WORKDIR /tmp
RUN VERSION=v2.3.1 && \
curl -sSL -o argocd https://github.com/argoproj/argo-cd/releases/download/$VERSION/argocd-linux-amd64 && \
chmod +x argocd

FROM golang:1.18-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o plugin cmd/argocd-plugin/main.go

FROM alpine:3.15
WORKDIR /app
USER 1000
COPY --from=downloader /tmp/argocd /usr/local/bin/argocd
COPY --from=build /app/plugin ./plugin
#RUN adduser -D -S argo -u 8737 

ENTRYPOINT [ "/app/plugin" ]