FROM curlimages/curl:7.81.0 as downloader
USER root
RUN VERSION=v2.2.5 && \
  curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/download/$VERSION/argocd-linux-amd64 && \
  chmod +x /usr/local/bin/argocd

FROM python:3.8-alpine3.15

WORKDIR /tmp
COPY --from=downloader /usr/local/bin/argocd /usr/local/bin/argocd
COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt
