package restapi

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"addreality_t/internal/addreality/users"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

// RESTAPI represents a REST API business logic server.
type RESTAPI struct {
	server http.Server
	errors chan error
	logger *zap.SugaredLogger
}

// New returns a new instance of the REST API server.
func New(logger *zap.SugaredLogger, port int, userController *users.Controller) *RESTAPI {
	router := mux.NewRouter()

	service := ServiceHandlers{
		UserController: userController,
	}
	router.HandleFunc("/", service.CountUser()).Methods(http.MethodGet)
	router.HandleFunc("/count", service.GetRobotUser()).Methods(http.MethodGet)

	return &RESTAPI{
		server: http.Server{
			Addr:    net.JoinHostPort("", strconv.Itoa(port)),
			Handler: router,
		},
		errors: make(chan error, 1),
		logger: logger,
	}
}

// Start diagnostics server.
func (rapi *RESTAPI) Start() {
	go func() {
		rapi.errors <- rapi.server.ListenAndServe()
		close(rapi.errors)
	}()
}

// Stop diagnostics server.
func (rapi *RESTAPI) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return rapi.server.Shutdown(ctx)
}

// Notify returns a channel to notify the caller about errors.
// If you receive an error from the channel you should stop the application.
func (rapi *RESTAPI) Notify() <-chan error {
	return rapi.errors
}
