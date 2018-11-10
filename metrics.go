package teleinfo

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	tagErrorType = "error_type"
)

var (
	frameReadCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "teleinfo_frames_read_total",
		Help: "The total number of read Teleinfo frames",
	})
	frameReadErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "teleinfo_frames_read_errors_total",
			Help: "The total number of frame read errors",
		},
		[]string{tagErrorType},
	)

	frameDecodedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "teleinfo_frames_decoded_total",
		Help: "The total number of decoded frames",
	})
	frameDecodeErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "teleinfo_frames_decode_errors_total",
			Help: "The total number of frame decoding errors",
		},
		[]string{tagErrorType},
	)
)

func incrementErrorCounter(counter *prometheus.CounterVec, errorType string) {
	counter.With(prometheus.Labels{tagErrorType: errorType}).Inc()
}