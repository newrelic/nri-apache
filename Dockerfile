FROM golang:1.10 as builder-apache
RUN go get -d github.com/newrelic/nri-apache/... && \
    cd /go/src/github.com/newrelic/nri-apache && \
    make && \
    strip ./bin/nr-apache

FROM newrelic/infrastructure:latest
COPY --from=builder-apache /go/src/github.com/newrelic/nri-apache/bin/nr-apache /var/db/newrelic-infra/newrelic-integrations/bin/nr-apache
COPY --from=builder-apache /go/src/github.com/newrelic/nri-apache/apache-definition.yml /var/db/newrelic-infra/newrelic-integrations/definition.yaml
