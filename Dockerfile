FROM golang:1.14 as builder

RUN git clone \
    --depth 1 \
    --single-branch \
    --branch=dev.go2go \
    --progress  \
    https://go.googlesource.com/go /go2

ENV CGO_ENABLED=0
WORKDIR /go2/src
RUN ./all.bash

FROM alpine:3.12.0

COPY --from=builder /go2 /go
ENV PATH=$PATH:/go/bin
ENV GOROOT=/go