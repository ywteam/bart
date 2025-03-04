package metrics_test

import (
	"strings"
	"testing"
	metrics "yellowteam/lib/metrics"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestHttpRequestCountWithPath(t *testing.T) {
	metrics.HttpRequestCountWithPath.WithLabelValues("/test").Inc()
	metrics.HttpRequestCountWithPath.WithLabelValues("/test").Inc()

	expected := `
		# HELP http_requests_total_with_path Number of HTTP requests by path.
		# TYPE http_requests_total_with_path counter
		http_requests_total_with_path{url="/test"} 2
	`
	if err := testutil.CollectAndCompare(metrics.HttpRequestCountWithPath, strings.NewReader(expected)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestHttpRequestDuration(t *testing.T) {
	metrics.HttpRequestDuration.WithLabelValues("/test").Observe(1.2)
	metrics.HttpRequestDuration.WithLabelValues("/test").Observe(0.8)

	expected := `
		# HELP http_request_duration_seconds Response time of HTTP request.
		# TYPE http_request_duration_seconds histogram
		http_request_duration_seconds_bucket{path="/test",le="0.005"} 0
		http_request_duration_seconds_bucket{path="/test",le="0.01"} 0
		http_request_duration_seconds_bucket{path="/test",le="0.025"} 0
		http_request_duration_seconds_bucket{path="/test",le="0.05"} 0
		http_request_duration_seconds_bucket{path="/test",le="0.1"} 0
		http_request_duration_seconds_bucket{path="/test",le="0.25"} 0
		http_request_duration_seconds_bucket{path="/test",le="0.5"} 0
		http_request_duration_seconds_bucket{path="/test",le="1"} 1
		http_request_duration_seconds_bucket{path="/test",le="2.5"} 2
		http_request_duration_seconds_bucket{path="/test",le="5"} 2
		http_request_duration_seconds_bucket{path="/test",le="10"} 2
		http_request_duration_seconds_sum{path="/test"} 2
		http_request_duration_seconds_count{path="/test"} 2
	`
	if err := testutil.CollectAndCompare(metrics.HttpRequestDuration, strings.NewReader(expected)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}
