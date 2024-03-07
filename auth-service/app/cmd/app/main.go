package main

import (
	"context"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/api/emailapi"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/handlers"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/repositories/postgres"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/repositories/sessionsrepo"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/repositories/usersrepo"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/servers"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/services/sessionsservice"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/services/usersservice"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/encryptor"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/logger"
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

	crypt := encryptor.NewEncryptor([]byte(getenv("AUTH_CRYPTO_KEY")))

	db, err := postgres.NewPostgresDB(&postgres.DBConfigs{
		Host:     getenv("DB_HOST"),
		Port:     getenv("DB_PORT"),
		Username: getenv("DB_USERNAME"),
		DBName:   getenv("DB_NAME"),
		Password: getenv("DB_PASSWORD"),
	})
	if err != nil {
		return err
	}

	emailAPI := emailapi.NewEmailAPI(ctx, getenv("EMAIL_SERVICE_HOST"))

	userRepository, err := usersrepo.NewUsersRepository(db)
	if err != nil {
		return err
	}

	userService := usersservice.NewUsersService(
		ctx,
		userRepository,
		emailAPI,
		crypt,
	)

	sessionsRepository, err := sessionsrepo.NewSessionsRepository(db)
	if err != nil {
		return err
	}

	ttl, err := strconv.Atoi(getenv("SESSION_EXPIRATION_TIME_MS"))
	if err != nil {
		return err
	}

	sessionTTL := time.Duration(ttl) * time.Millisecond

	sessionsService := sessionsservice.NewSessionsService(ctx, sessionTTL, sessionsRepository)

	handler := handlers.NewHTTPHandler(
		userService,
		sessionsService,
		sessionTTL,
	)

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
