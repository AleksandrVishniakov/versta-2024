package main

import (
	"context"
	"errors"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/api/authapi"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/apiclient"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/handlers"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/repositories/chattersrepo"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/repositories/messagesrepo"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/repositories/postgres"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/servers"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/chatters"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/chattokens"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/hubs"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/messages"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/jwttokens"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/logger"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/scrambler"
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

	aesScrambler := scrambler.NewAES256([]byte(getenv("CHAT_CRYPTO_KEY")))

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

	messagesRepo := messagesrepo.NewMessagesRepository(db)
	messagesStorage := messages.NewMessagesStorage(messagesRepo, aesScrambler)

	chattersRepo := chattersrepo.NewChattersRepository(db)
	chattersStorage := chatters.NewChattersStorage(chattersRepo)

	hubManager := hubs.NewHubManager(
		ctx,
		messagesStorage,
	)

	accessTokenTTL, err := strconv.Atoi(getenv("ACCESS_TOKEN_TTL_MS"))
	if err != nil {
		return err
	}

	refreshTokenTTL, err := strconv.Atoi(getenv("REFRESH_TOKEN_TTL_MS"))
	if err != nil {
		return err
	}

	tokensManager := jwttokens.NewTokensManager(
		[]byte(getenv("JWT_SIGNATURE_KEY")),
		time.Duration(accessTokenTTL)*time.Millisecond,
		time.Duration(refreshTokenTTL)*time.Millisecond,
	)

	authAPI := authapi.NewAuthAPI(ctx, getenv("AUTH_SERVICE_HOST"), apiclient.NewAPIClient(ctx))
	admin, err := authAPI.Admin()
	if err != nil {
		return err
	}

	chatter, err := chattersStorage.FindByUserId(admin.Id)
	var adminID = 0
	if err != nil && !errors.Is(err, chatters.ErrChatterNotFound) {
		return err
	}

	if chatter != nil {
		adminID = chatter.Id
	} else {
		id, err := chattersStorage.CreateWithId(admin.Id)
		if err != nil {
			return err
		}
		adminID = id
	}

	handler := handlers.NewHTTPHandler(
		ctx,
		hubManager,
		tokensManager,
		chattersStorage,
		messagesStorage,
		chattokens.NewChatTokens(),
		adminID,
	)

	server := servers.NewHTTPServer(httpPort, handler.Handler())

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
