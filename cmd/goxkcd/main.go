package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	httpServer "github.com/toadharvard/goxkcd/internal/api/http"
	"github.com/toadharvard/goxkcd/internal/app"
	"github.com/toadharvard/goxkcd/internal/config"
)

func getValuesFromArgs() (string, string, int) {
	configPath := flag.String("c", config.DefaultConfigPath, "Config path")
	host := flag.String("h", "localhost", "Host")
	port := flag.Int("p", 8080, "Port")
	flag.Parse()
	return *configPath, *host, *port
}

func run() error {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)

	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: lvl,
			},
		),
	)

	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	configPath, hostArg, portArg := getValuesFromArgs()
	cfg, err := config.New(configPath)
	if err != nil {
		return err
	}

	app, err := app.New(cfg)
	if err != nil {
		return err
	}

	host := cfg.HTTPServer.Host
	if hostArg != "" {
		host = hostArg
	}

	port := cfg.HTTPServer.Port
	if portArg != 0 {
		port = portArg
	}

	err = httpServer.Run(ctx, app, host, port, cfg.ComixUpdateInterval)
	return err
}

func main() {
	if err := run(); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
