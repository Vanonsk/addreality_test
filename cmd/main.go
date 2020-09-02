package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	usercontr "addreality_t/internal/addreality/users"
	userdb "addreality_t/internal/addreality/users/withdb"
	"addreality_t/internal/metrics"
	"addreality_t/internal/resources"
	"addreality_t/internal/restapi"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Error when close logger")
		}
	}()

	slogger := logger.Sugar()
	slogger.Info("Starting the application...")
	slogger.Info("Reading configuration and initializing resources...")

	rsc, err := resources.New(slogger)
	if err != nil {
		slogger.Fatalw("Can't initialize resources.", "err", err)
	}
	defer func() {
		rsc.Release()
	}()

	slogger.Info("Configuring the application units...")
	userdbLevel := userdb.New(rsc.DB)
	uc := usercontr.NewController(userdbLevel)

	slogger.Info("Starting the servers...")
	rapi := restapi.New(slogger, rsc.Config.RESTAPIPort, uc)
	rapi.Start()

	metricsServer := metrics.New(slogger, rsc.Config.MetricsPort)
	metricsServer.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case x := <-interrupt:
		slogger.Infow("Received a signal.", "signal", x.String())
	case err := <-metricsServer.Notify():
		slogger.Errorw("Received an error from the metrics server.", "err", err)
	case err := <-rapi.Notify():
		slogger.Errorw("Received an error from the business logic server.", "err", err)
	}

	slogger.Info("Stopping the servers...")
	err = rapi.Stop()
	if err != nil {
		slogger.Error("Got an error while stopping the business logic server.", "err", err)
	}

	err = metricsServer.Stop()
	if err != nil {
		slogger.Error("Got an error while stopping the metrics logic server.", "err", err)
	}

	slogger.Info("The app is calling the last defers and will be stopped.")
}
