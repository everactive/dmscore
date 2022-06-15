FROM golang:1.16 as builder1
COPY . /iot-identity
WORKDIR /iot-identity
ARG GIT_SUMMARY
RUN CGO_ENABLED=1 GOOS=linux go build -a -o /go/bin/identity -ldflags='-extldflags "-static" -X github.com/everactive/dmscore/iot-identity/version.Version='"$GIT_SUMMARY" bin/identity/identity.go

# Copy the built applications to the docker image
FROM ubuntu:18.04
WORKDIR /srv
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY --from=builder1 /go/bin/identity /srv/identity

EXPOSE 8030
ENTRYPOINT /srv/identity run
