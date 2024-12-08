package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "education-api",
	Short: "Education API - A RESTful service for managing educational resources",
	Long: `A RESTful API service built with Go that provides endpoints 
for managing educational resources, courses, and related data.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
