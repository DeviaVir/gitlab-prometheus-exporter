FROM golang:alpine AS builder

ENV GO111MODULE="on"
ENV CGO_ENABLED="0"

RUN apk add --update git

RUN mkdir -p /go/src/github.com/DeviaVir/gitlab-prometheus-exporter

COPY . /go/src/github.com/DeviaVir/gitlab-prometheus-exporter

RUN cd /go/src/github.com/DeviaVir/gitlab-prometheus-exporter \
 && go mod vendor \
 && go build \
      -mod vendor \
      -o /go/bin/gitlab-prometheus-exporter

FROM alpine
COPY --from=builder /go/bin/gitlab-prometheus-exporter /usr/local/bin/gitlab-prometheus-exporter
CMD ["/usr/local/bin/gitlab-prometheus-exporter"]
