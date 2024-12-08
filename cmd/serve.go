package cmd

import (
	"education-api/server"
	"fmt"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Education API server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Education API server...")
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
