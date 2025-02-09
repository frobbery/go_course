package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	//nolint:depguard
	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/app"
	//nolint:depguard
	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/logger"
	//nolint:depguard
	internalhttp "github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/server/http"
	//nolint:depguard
	memorystorage "github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/storage/memory"
	//nolint:depguard
	sqlstorage "github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := Read(configFile)
	if err != nil {
		panic("could not init config: " + err.Error())
	}

	logger := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := initStorage(ctx, config.DB)
	if err != nil {
		panic("could not init config: " + err.Error())
	}
	defer func() {
		if err := storage.Close(); err != nil {
			logger.Error("cannot close psql connection: " + err.Error())
		}
	}()

	calendar := app.New(logger, storage)

	server := internalhttp.NewServer(logger, calendar, config.HTTP.Host, config.HTTP.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func initStorage(ctx context.Context, dbConfig DBConfig) (storage app.Storage, err error) {
	if dbConfig.InMemory {
		storage = memorystorage.New()
	} else {
		storage = sqlstorage.New()
	}
	if err := storage.Connect(ctx, dbConfig.DSN); err != nil {
		return nil, fmt.Errorf("cannot connect to psql: %w", err)
	}
	if err := storage.Migrate(ctx, dbConfig.Migration); err != nil {
		return nil, fmt.Errorf("cannot migrate: %w", err)
	}
	return storage, nil
}
