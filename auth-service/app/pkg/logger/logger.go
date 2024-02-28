package logger

import (
	"io"
	"log"
	"log/slog"
)

const (
	LevelProduction = slog.Level(-2)
)

var levelNames = map[slog.Leveler]string{
	LevelProduction: "PRODUCTION",
}

func InitLogger(writer io.Writer, logLevel string) {
	level := parseLogLevel(logLevel)

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
		AddSource: level != LevelProduction,
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
		return LevelProduction
	case "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		log.Printf("logger: cannot define log level %s. log level is set as INFO", levelName)
		return slog.LevelInfo
	}
}
