package main

import (
	"io"
	"log"
	"os"

	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/logger"
)

const logFileName = "app.log"

func main() {
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("%s file opening error: %s", logFileName, err.Error())
	}
	defer logFile.Close()

	logger.InitLogger(io.MultiWriter(os.Stdout, logFile), os.Getenv)
}
