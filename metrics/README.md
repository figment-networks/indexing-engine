# Metrics

In `metrics` package you can find simple implementation of metrics engine, that is metric system agnostic.
Currently available engines are:
  * Prometheus (https://github.com/prometheus/client_golang)

Nice to have in future:
  * Statsd (eg. https://github.com/alexcesaro/statsd)
  * Expvar (https://golang.org/pkg/expvar/)

The example of how to write an engine lives in `/metrics/metrics_test.go`
