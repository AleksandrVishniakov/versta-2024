package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/AleksandrVishniakov/versta-2024/email-service/app/internal/handlers"
	"github.com/AleksandrVishniakov/versta-2024/email-service/app/internal/servers"
	"github.com/AleksandrVishniakov/versta-2024/email-service/app/internal/services/emailservice"
	"github.com/AleksandrVishniakov/versta-2024/email-service/app/pkg/logger"
)

const logFileName = "logs/app.log"

func run(
	ctx context.Context,
	getenv func(string) string,
	writer io.Writer,
	errorsCh chan<- error,
) error {
	var (
		httpPort = getenv("HTTP_PORT")
		logLevel = getenv("LOG_LEVEL")
	)

	logger.InitLogger(writer, logLevel)

	emailService := emailservice.NewEmailService(&emailservice.EmailConfigs{
		Host:        getenv("EMAIL_HOST"),
		Port:        getenv("EMAIL_PORT"),
		SenderEmail: getenv("EMAIL_SENDER"),
		Password:    getenv("EMAIL_PASSWORD"),
	})

	handler := handlers.NewHTTPHandler(emailService)

	server := servers.NewHTTPServer(ctx, httpPort, handler.Handler())

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
		errorsCh <- err
		return nil
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

	var errorsCh = make(chan error)
	err = run(ctx, os.Getenv, logWriter, errorsCh)
	if err != nil {
		slog.Error("server starting error", slog.String("error", err.Error()))
	}

	if err := <-errorsCh; err != nil {
		slog.Error("server running error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
