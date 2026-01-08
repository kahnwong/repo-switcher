package main

import (
	"os"

	"github.com/kahnwong/repo-switcher/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
