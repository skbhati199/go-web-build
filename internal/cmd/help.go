package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:     "help [command]",
	Aliases: []string{"h", "?", "manual"},
	Short:   "Show help for any command",
	Long: `Display detailed help and usage information for any command.
Includes examples, available flags, and aliases.`,
	RunE: showHelp,
}

func showHelp(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return rootCmd.Help()
	}

	targetCmd, _, err := rootCmd.Find(args)
	if err != nil {
		return fmt.Errorf("unknown command: %s", strings.Join(args, " "))
	}

	return targetCmd.Help()
}

func init() {
	rootCmd.AddCommand(helpCmd)
}
