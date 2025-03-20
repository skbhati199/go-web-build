package cmd

import (
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:     "dev",
	Aliases: []string{"d", "serve", "start", "watch", "run"},
	Short:   "Start development server",
	Long: `Start a development server with hot reload and debugging capabilities.

Development Features:
  • Hot Module Replacement (HMR)
  • Source maps support
  • Fast refresh
  • Development tools integration
  • Live reload
  
Server Options:
  --port          Specify server port (default: 3000)
  --host          Specify host address (default: localhost)
  --https         Enable HTTPS
  --proxy         Configure API proxy
  --open          Open in default browser`,
	Example: `  # Start server with default settings
  gobuild dev

  # Start on custom port with HTTPS
  gobuild dev --port 8080 --https

  # Start with API proxy
  gobuild dev --proxy /api:http://localhost:8000`,
}

func init() {
	devCmd.Flags().IntP("port", "p", 3000, "development server port")
	devCmd.Flags().StringP("host", "H", "localhost", "host address")
	devCmd.Flags().BoolP("https", "S", false, "enable HTTPS")
	devCmd.Flags().StringP("proxy", "P", "", "API proxy configuration")
	devCmd.Flags().BoolP("open", "o", false, "open in browser")
	devCmd.Flags().BoolP("hot", "h", true, "enable hot reload")

	rootCmd.AddCommand(devCmd)
}
