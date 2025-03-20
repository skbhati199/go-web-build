package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b", "make", "compile", "bundle"},
	Short:   "Build the web application",
	Long: `Build your web application for development or production.
Supports various optimization options and build modes.`,
	Example: `  # Build for production
  gobuild build --mode production

  # Build with source maps
  gobuild build --mode development --sourcemap

  # Build with custom output directory
  gobuild build --out-dir ./dist`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mode, _ := cmd.Flags().GetString("mode")

		fmt.Printf("Building project in %s mode\n", mode)
		return nil
	},
}

func init() {
	buildCmd.Flags().StringP("mode", "m", "production", "build mode (development, production)")
	buildCmd.Flags().StringP("out", "o", "dist", "output directory")
	buildCmd.Flags().BoolP("minify", "M", true, "enable minification")
	buildCmd.Flags().BoolP("sourcemap", "s", false, "generate source maps")

	rootCmd.AddCommand(buildCmd)
}
