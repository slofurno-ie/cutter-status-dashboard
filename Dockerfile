## Local development build
FROM golang:1.13-buster AS dev_environment

WORKDIR /app
RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]

COPY . .

# meant to be used with docker-compose and volume mounted ssh keys.
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
ENTRYPOINT CompileDaemon -log-prefix=false -build="go build ." -command "./cutter-status-dashboard"

FROM golang:1.14.4-alpine3.12 AS builder

ARG SSH_KEY
WORKDIR /app
COPY . .

RUN apk add --no-cache openssh-client git \
  && mkdir ~/.ssh \
  && ssh-keyscan github.com > ~/.ssh/known_hosts

RUN echo "${SSH_KEY}" > /root/.ssh/id_rsa && \
    chmod 600 /root/.ssh/id_rsa && \
    git config --global url."git@github.com:".insteadOf "https://github.com/" && \
    ssh-keyscan -t rsa github.com > /root/.ssh/known_hosts && \
    GOPRIVATE=github.com/IdeaEvolver/ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o cutter-status-dashboard

EXPOSE 8080