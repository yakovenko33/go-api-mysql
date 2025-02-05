package cmd

import (
	"fmt"
	"os"

	createS_uper_admin "go-api-docker/cmd/commands/CreateSuperAdmin"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-cli-app",
	Short: "A simple CLI app similar to Laravel Artisan",
	Long:  `This CLI app is built using Cobra in Go, just like Artisan commands in Laravel.`,
}

func init() {
	// Add the custom greet command to the root command
	rootCmd.AddCommand(createS_uper_admin.СreateSuperAdmin)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
