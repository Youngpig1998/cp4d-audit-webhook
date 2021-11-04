FROM registry.access.redhat.com/ubi8/ubi:latest AS builder

LABEL maintainer="Watson Platform" \
      version="1.0.0"

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPATH="/go" 

ENV PATH=$PATH:/usr/local/go/bin

RUN yum install -y wget git tar gzip \
  && wget https://golang.org/dl/go1.14.6.linux-amd64.tar.gz \
  && tar -xvf go1.14.6.linux-amd64.tar.gz -C /usr/local \
  && mkdir -p /go/bin /go/src /go/pkg \
  && yum clean all

COPY audit-webhook /go/src/audit-webhook

WORKDIR /go/src/audit-webhook

RUN go build

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

COPY --from=builder /go/src/audit-webhook/audit-webhook /audit-webhook

RUN microdnf install -y shadow-utils \
    && adduser audit -u 10001 -g 0 \
    && chown audit:root /audit-webhook \
    && chmod +x /audit-webhook 

USER 10001

ENTRYPOINT ["/audit-webhook"]