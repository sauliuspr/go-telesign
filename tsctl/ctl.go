package main

import (
	"github.com/sauliuspr/go-telesign/tsctl/command"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:        "tsctl",
		Short:      "A command line client for Telesign",
		SuggestFor: []string{"tsctl"},
	}
)

func main() {

	rootCmd.AddCommand(
		command.NewStatusCommand(),
	)

}
