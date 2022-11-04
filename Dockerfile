# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN apk add build-base
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w" -o plugin cmd/argocd-plugin/main.go

FROM alpine:3.16.2

RUN apk add diffutils

USER 1000

COPY --from=build /app/plugin /

CMD [ "/plugin" ]
