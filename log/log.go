package log

import (
	"github.com/rs/zerolog"
	"io"
	"os"
)

var logger zerolog.Logger

func InitLogger(path string) {
	var w io.Writer
	switch path {
	case "stderr":
		w = os.Stderr
	default:
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		w = f
	}
	logger = zerolog.New(w).With().Timestamp().Logger()
}

func GetLogger(moduleName string) zerolog.Logger {
	return logger.With().Str("from", moduleName).Logger()
}

func GenCheck(l zerolog.Logger) func(err error) {
	return func(err error) {
		if err != nil {
			l.Panic().Msg(err.Error())
		}
	}
}
