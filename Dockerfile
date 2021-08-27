FROM golang:1.10 as builder-apache
COPY . /go/src/github.com/newrelic/nri-apache/
RUN cd /go/src/github.com/newrelic/nri-apache && \
    make && \
    strip ./bin/nri-apache

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder-apache /go/src/github.com/newrelic/nri-apache/bin/nri-apache /nri-sidecar/newrelic-infra/newrelic-integrations/bin/nri-apache
USER 1000
