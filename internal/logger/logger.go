package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func Init(debug bool) {
	zerolog.TimeFieldFormat = time.RFC3339

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}

	log.Logger = zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

func Fatal(err error, msg string) {
	log.Fatal().Err(err).Msg(msg)
}
