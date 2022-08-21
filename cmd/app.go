package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/sn5ake6/otus-final-project/internal/app"
	"github.com/sn5ake6/otus-final-project/internal/bucket"
	"github.com/sn5ake6/otus-final-project/internal/config"
	"github.com/sn5ake6/otus-final-project/internal/logger"
	internalgrpc "github.com/sn5ake6/otus-final-project/internal/server/grpc"
	internalhttp "github.com/sn5ake6/otus-final-project/internal/server/http"
	"github.com/sn5ake6/otus-final-project/internal/version"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		version.PrintVersion()
		return
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.NewConfig()
	err := config.LoadConfig(configFile, &cfg)
	if err != nil {
		return err
	}

	logg, err := logger.New(cfg.Logger.Level)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := app.NewStorage(cfg.Storage)
	err = storage.Connect(ctx)
	if err != nil {
		return err
	}

	bucket := bucket.NewLeakyBucket(ctx, cfg.Limit)

	service := app.New(logg, storage, bucket)

	startGRPCServer(ctx, cfg, logg, service)

	startHTTPServer(ctx, cfg, logg, service)

	logg.Info("anti-brute-force is running...")

	<-ctx.Done()

	return nil
}

func startGRPCServer(ctx context.Context, cfg config.Config, logg *logger.Logger, service *app.App) {
	grpcAddr := net.JoinHostPort(cfg.GRPCServer.Host, cfg.GRPCServer.Port)
	grpcServer := internalgrpc.NewServer(grpcAddr, logg, service)

	go func() {
		if err := grpcServer.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()
}

func startHTTPServer(ctx context.Context, cfg config.Config, logg *logger.Logger, service *app.App) {
	httpAddr := net.JoinHostPort(cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	httpServer := internalhttp.NewServer(httpAddr, logg, service)

	go func() {
		if err := httpServer.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()
}
