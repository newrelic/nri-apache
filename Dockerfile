FROM golang:1.10 as builder-apache
RUN go get -d github.com/newrelic/nri-apache/... && \
    cd /go/src/github.com/newrelic/nri-apache && \
    make && \
    strip ./bin/nr-apache

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder-apache /go/src/github.com/newrelic/nri-apache/bin/nr-apache /nri-sidecar/newrelic-infra/newrelic-integrations/bin/nr-apache
COPY --from=builder-apache /go/src/github.com/newrelic/nri-apache/apache-definition.yml /nri-sidecar/newrelic-infra/newrelic-integrations/definition.yaml
USER 1000
