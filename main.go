package main

import (
	"github.com/spf13/cobra"
	"os"
	"outback/kingo/service/crawl/cmd"
	"outback/kingo/service/flag"
)

func main() {
	rootCmd := NewCmdRoot()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	// time.Sleep(time.Hour)
}

func NewCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "prepare <subcommand> [flags]",
		Short: "prepare agentless scan",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	flag.AddOption(rootCmd)
	rootCmd.AddCommand(cmd.NewCrawlCommand())
	rootCmd.AddCommand(cmd.NewCrawlNameCommand())

	return rootCmd
}
