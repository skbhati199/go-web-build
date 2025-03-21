package cmd

import (
	"github.com/skbhati199/go-web-build/internal/config"
	"github.com/skbhati199/go-web-build/internal/recovery"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
	debug   bool
	env     string
)

var rootCmd = &cobra.Command{
	Use:     "gobuild",
	Aliases: []string{"gb", "gwb", "webuild"},
	Short:   "A modern web application build tool",
	Long: `GoBuild is a powerful CLI tool for creating and building modern web applications.

Key Features:
  • Project scaffolding with multiple framework support
  • Optimized build configurations for development and production
  • Development server with hot reload capability
  • Built-in template system for quick project setup
  • Extensible plugin architecture

For detailed documentation, visit: https://github.com/skbhati199/go-web-build`,
	Example: `  # Create a new React project
  gobuild create myapp --framework react --template typescript

  # Start development server with hot reload
  gobuild dev --port 3000 --hot

  # Build for production with optimizations
  gobuild build --mode production --minify

  # Deploy your application
  gobuild deploy --target aws`,
}

func Execute() error {
	recovery := recovery.NewRecoveryHandler(debug)

	return recovery.WrapHandler(func() error {
		return rootCmd.Execute()
	})()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug mode")
	rootCmd.PersistentFlags().StringVar(&env, "env", "", "environment (development, staging, production)")

	// Remove duplicate flag definitions
	// rootCmd.PersistentFlags().StringP("config", "c", "", "config file path")
	// rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug mode")
	// rootCmd.PersistentFlags().StringP("env", "e", "", "environment (development, staging, production)")
}
