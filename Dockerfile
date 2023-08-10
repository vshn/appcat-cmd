FROM docker.io/library/golang:alpine AS build

WORKDIR /src
ENV CGO_ENABLED=0

COPY . .
RUN go build -v -o /appcat-cli .

# Runtime
FROM docker.io/library/alpine:latest

ENTRYPOINT ["/bin/appcat-cli"]
COPY --from=build /appcat-cli /bin/appcat-cli
