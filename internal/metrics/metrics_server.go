package metrics

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Prometheus metrics server
type MetricsAPI struct {
	server http.Server
	errors chan error
	logger *zap.SugaredLogger
}

// New. Return init metrics server
func New(logger *zap.SugaredLogger, port int) *MetricsAPI {
	http.Handle("/metrics", promhttp.Handler())
	return &MetricsAPI{
		server: http.Server{
			Addr:    net.JoinHostPort("", strconv.Itoa(port)),
			Handler: nil,
		},
		errors: make(chan error, 1),
		logger: logger,
	}
}

// Start metrics server.
func (d *MetricsAPI) Start() {
	go func() {
		d.logger.Info("Starting prometheus metrics...")
		d.errors <- d.server.ListenAndServe()
		close(d.errors)
	}()
}

// Stop metrics server.
func (d *MetricsAPI) Stop() error {
	d.logger.Info("Stop prometheus metrics...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.server.Shutdown(ctx)
}

// Notify returns a channel to notify the caller about errors.
// If you receive an error from the channel metrics you should stop the application.
func (d *MetricsAPI) Notify() <-chan error {
	return d.errors
}
