package logger

import (
	"github.com/rs/zerolog/log"
)

func Info(message string) {
	log.Info().Msg(message)
}

func Error(message string) {
	log.Error().Msg(message)
}

func Fatal(message string) {
	log.Fatal().Msg(message)
}

func Debug(message string) {
	log.Debug().Msg(message)
}
