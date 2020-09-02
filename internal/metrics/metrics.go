package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	UserRequestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_request_total",
		Help: "The total number of request",
	})
	UserStatusRequestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_user_status_request_total",
		Help: "The total number of user status request",
	})
)
