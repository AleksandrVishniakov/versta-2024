package main

import (
	"context"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/handlers"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/repositories/postgres"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/servers"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/logger"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
)

const logFileName = "logs/app.log"

func run(
	ctx context.Context,
	getenv func(string) string,
	writer io.Writer,
) error {
	var (
		httpPort = getenv("HTTP_PORT")
		logLevel = getenv("LOG_LEVEL")
	)

	logger.InitLogger(writer, logLevel)

	_, err := postgres.NewPostgresDB(&postgres.DBConfigs{
		Host:     getenv("DB_HOST"),
		Port:     getenv("DB_PORT"),
		Username: getenv("DB_USERNAME"),
		DBName:   getenv("DB_NAME"),
		Password: getenv("DB_PASSWORD"),
	})

	if err != nil {
		return err
	}

	handler := handlers.NewHTTPHandler()

	server := servers.NewHTTPServer(httpPort, handler)

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("server is shutting down")

				if err := server.Shutdown(ctx); err != nil {
					slog.Error("error shutting down http server", slog.String("error", err.Error()))
				}
			}

		}
	}()

	slog.Info("server started on http://localhost:" + httpPort)
	if err := server.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("%s file opening error: %s", logFileName, err.Error())
	}
	defer logFile.Close()

	logWriter := io.MultiWriter(os.Stdout, logFile)

	if err := run(ctx, os.Getenv, logWriter); err != nil {
		slog.Error("server running error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
