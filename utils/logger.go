package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/Ether-Security/leviathan/libs"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLog(options *libs.Options) {
	var writers []io.Writer

	// Check if this is a JSON output or a pretty one
	if options.Log.JSON {
		writers = append(writers, os.Stdout)
	} else {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Add file writers
	if _, err := os.Stat(options.Log.Directory); os.IsNotExist(err) {
		os.Mkdir(options.Log.Directory, 0700)
	}
	tempFile, err := os.CreateTemp(options.Log.Directory, fmt.Sprintf("%s-*.log", libs.NAME))
	if err == nil {
		writers = append(writers, tempFile)
	}

	mw := io.MultiWriter(writers...)
	Logger = zerolog.New(mw).With().Timestamp().Logger()

	if options.Log.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if options.Log.Quiet {
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	Logger.Info().Msgf("Logs are stored in %s", tempFile.Name())
	Logger.Info().Msgf("Logs level:  %s", Logger.GetLevel().String())
}
