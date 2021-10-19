package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/eldeal/observation-controller/config"
	"github.com/eldeal/observation-controller/service"
)

func main() {
	log.Namespace = "observation-controller"
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(ctx, "application unexpectedly failed", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run(ctx context.Context) error {
	svcList := service.NewServiceList(&service.Init{})
	svcErrors := make(chan error, 1)

	cfg, err := config.Get()
	if err != nil {
		log.Error(ctx, "unable to retrieve service configuration", err)
		return err
	}

	log.Info(ctx, "got service configuration", log.Data{"config": cfg})

	svc := service.New()
	if err := svc.Init(ctx, cfg, svcList); err != nil {
		log.Error(ctx, "failed to initialise service", err)
		return err
	}
	svc.Run(ctx, svcErrors)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-svcErrors:
		log.Error(ctx, "service error received", err)
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}

	return svc.Close(ctx)
}
