# Build
FROM golang:latest AS build

ENV HOME=/appcat-cli

WORKDIR ${HOME}

COPY . ${HOME}

RUN go build -v .

# Runtime
FROM quay.io/vshn/k8ify:latest

COPY --from=build appcat-cli /bin/appcat
