package main

import (
	"log/slog"
	"os"

	"github.com/kahnwong/repo-switcher/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		slog.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}
