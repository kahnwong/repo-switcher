package main

import (
	"os"

	"github.com/kahnwong/repo-switcher/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// Set default log level to info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("command execution failed")
	}
}
