FROM golang:1.17 as builder1
COPY . /iot-devicetwin
WORKDIR /iot-devicetwin
ARG GIT_SUMMARY
RUN CGO_ENABLED=1 GOOS=linux go build -a -o /go/bin/devicetwin -ldflags='-extldflags "-static" -X github.com/everactive/dmscore/iot-devicetwin/version.Version='"$GIT_SUMMARY" bin/devicetwin/devicetwin.go

# Copy the built applications to the docker image
FROM ubuntu:18.04
WORKDIR /srv
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY --from=builder1 /go/bin/devicetwin /srv/devicetwin

# Note: /srv/certs is handled by the k8s file


EXPOSE $PORT
ENTRYPOINT /srv/devicetwin run