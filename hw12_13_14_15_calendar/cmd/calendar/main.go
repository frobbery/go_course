package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/app"
	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/storage/memory"
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

	logg := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := initStorage(ctx, config.DB)
	if err != nil {
		panic("could not init config: " + err.Error())
	}
	defer func() {
		if err := storage.Close(); err != nil {
			logg.Error("cannot close psql connection", err)
		}
	}()

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func initStorage(ctx context.Context, dbConfig DbConfig) (storage app.Storage, err error) {
	if dbConfig.InMemory {
		storage = memorystorage.New()
	} else {
		storage = sqlstorage.New()
	}
	if err := storage.Connect(ctx, dbConfig.DSN); err != nil {
		return nil, fmt.Errorf("cannot connect to psql: %v", err)
	}
	if err := storage.Migrate(ctx, dbConfig.Migration); err != nil {
		return nil, fmt.Errorf("cannot migrate: %v", err)
	}
}
