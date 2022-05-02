FROM registry.access.redhat.com/ubi8/go-toolset:1.16.12-2 as builder
ENV GOPATH=/go/
USER root

WORKDIR /spi-file

# Copy the Go Modules manifests
COPY ./go.mod go.mod
COPY ./go.sum go.sum

# Copy the go sources
COPY ./server.go server.go
COPY ./websocket websocket
#COPY ./gitfile gitfile
COPY static/index.html static/index.html

# build service
# Note that we're not running the tests here. Our integration tests depend on a running cluster which would not be
# available in the docker build.
RUN export ARCH="$(uname -m)" && if [[ ${ARCH} == "x86_64" ]]; then export ARCH="amd64"; elif [[ ${ARCH} == "aarch64" ]]; then export ARCH="arm64"; fi && \
    CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -a -o spi-file server.go

FROM registry.access.redhat.com/ubi8-minimal:8.5

COPY --from=builder /spi-file/spi-file /spi-file
COPY --from=builder /spi-file/static/index.html /static/index.html

WORKDIR /
USER 65532:65532

ENTRYPOINT ["/spi-file"]

