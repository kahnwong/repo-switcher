package main

import (
	"github.com/kahnwong/repo-switcher/cmd"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("command execution failed")
	}
}
