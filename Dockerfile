FROM golang:1.18 as builder1
COPY . /dmscore
WORKDIR /dmscore
ARG GIT_SUMMARY
RUN CGO_ENABLED=1 GOOS=linux go build -a -o /go/bin/dmscore -ldflags='-extldflags "-static" -X github.com/everactive/dmscore/versions.Version='"$GIT_SUMMARY" bin/dmscore/dmscore.go

# Copy the built applications to the docker image
FROM ubuntu:18.04
WORKDIR /root/
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY --from=builder1 /go/bin/dmscore .

COPY db/migrations /migrations/dmscore
COPY iot-identity/db/migrations /migrations/identity
COPY iot-devicetwin/db/migrations /migrations/devicetwin

EXPOSE 8010
ENTRYPOINT /root/dmscore run
