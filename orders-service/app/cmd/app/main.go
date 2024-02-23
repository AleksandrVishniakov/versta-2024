package main

import (
	"context"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/servers"
	"io"
	"log"
	"os"

	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/logger"
)

const logFileName = "app.log"

var (
	httpPort = os.Getenv("HTTP_PORT")
	logLevel = os.Getenv("LOG_LEVEL")
)

func main() {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("%s file opening error: %s", logFileName, err.Error())
	}
	defer logFile.Close()

	logger.InitLogger(io.MultiWriter(os.Stdout, logFile), logLevel)

	server := servers.NewHTTPServer(httpPort, nil)

	if err := server.Run(); err != nil {
		err := server.Shutdown(ctx)
		if err != nil {
			log.Println("server shutdown error:", err.Error())
		}
	}
}
