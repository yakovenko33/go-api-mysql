package cli

import (
	"fmt"
	"os"

	creates_uper_admin "go-api-docker/cmd/cli/commands/create_super_admin"
	database "go-api-docker/internal/common/database"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-cli-app",
	Short: "A simple CLI app similar to Laravel Artisan",
	Long:  `This CLI app is built using Cobra in Go, just like Artisan commands in Laravel.`,
}

func init() {
	// Add the custom greet command to the root command

}

func Execute() {
	DB, err := database.ProvideDBConnection()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(creates_uper_admin.СreateSuperAdmin(DB))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
