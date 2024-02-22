package main

import (
	"io"
	"log"
	"log/slog"
)

const logLevelEnvKey = "LOG_LEVEL"

const (
	levelProduction = slog.Level(-2)
)

var levelNames = map[slog.Leveler]string{
	levelProduction: "PRODUCTION",
}

func initLogger(writer io.Writer, getenv func(string) string) {
	level := parseLogLevel(getenv(logLevelEnvKey))

	options := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := levelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
		AddSource: level != levelProduction,
	}

	logger := slog.New(slog.NewJSONHandler(writer, options))

	slog.SetDefault(logger)
}

func parseLogLevel(levelName string) slog.Level {
	switch levelName {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "PRODUCTION":
		return levelProduction
	case "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		log.Printf("logger: cannot define log level %s. log level is set as INFO", levelName)
		return slog.LevelInfo
	}
}
