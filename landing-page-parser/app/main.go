package main

import (
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const logFileName = "logs/app.log"

const (
	indexFilePath  = "web/app/build/index.html"
	staticFilesDir = "web/app/build"
)

func main() {
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("%s file opening error: %s", logFileName, err.Error())
	}
	defer logFile.Close()

	initLogger(io.MultiWriter(os.Stdout, logFile), os.Getenv)

	mux := http.NewServeMux()

	mux.Handle("/static/", http.FileServer(http.Dir(staticFilesDir)))
	mux.HandleFunc("/", Recovery(HTTPLogger(ErrorHandler(parsePage))))

	server := newHTTPServer(mux, os.Getenv)

	slog.Info("server started on http://localhost:" + os.Getenv("HTTP_PORT"))
	if err := server.run(); err != nil {
		log.Fatalf("server fatal error: %s", err.Error())
	}
}

func parsePage(w http.ResponseWriter, _ *http.Request) error {
	var pageConfigs = struct {
		OrdersServiceHost string
	}{
		OrdersServiceHost: os.Getenv("ORDERS_SERVICE_HOST"),
	}

	html, err := template.ParseFiles(indexFilePath)
	if err != nil {
		return err
	}

	err = html.Execute(w, pageConfigs)
	if err != nil {
		return err
	}

	return nil
}
