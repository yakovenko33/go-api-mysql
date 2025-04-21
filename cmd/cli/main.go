package main

import (
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	creates_uper_admin "go-api-docker/cmd/cli/commands/create_super_admin"
	fmt_color "go-api-docker/cmd/cli/commands/enums/fmt_color"
	database "go-api-docker/internal/common/database"
	access_control_model "go-api-docker/internal/common/security/access_control_models"
)

var rootCmd = &cobra.Command{
	Use:   "go-artisan",
	Short: "A simple CLI app similar to Laravel Artisan",
	Long:  `This CLI app is built using Cobra in Go, just like Artisan commands in Laravel.`,
}

func RunCLI() {
	DB, err := database.ProvideDBConnection()
	if err != nil {
		fmt_color.PrintError(err.Error())
		os.Exit(1)
	}
	accessControlModel, err := access_control_model.InitAccessControlModelForConsole(DB)
	if err != nil {
		fmt_color.PrintError(err.Error())
		os.Exit(1)
	}

	addCommand(DB, accessControlModel)
	if err := rootCmd.Execute(); err != nil {
		fmt_color.PrintError(err.Error())
		os.Exit(1)
	}
}

func addCommand(DB *gorm.DB, accessControlModel *casbin.Enforcer) {
	rootCmd.AddCommand(creates_uper_admin.СreateSuperAdmin(DB, accessControlModel))
}

func main() {
	RunCLI()
}
