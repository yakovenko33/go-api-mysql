package main

import (
	"go-api-docker/cmd"
	"go-api-docker/internal/common/database"
	logging "go-api-docker/internal/common/logging"
	"go-api-docker/internal/common/routes"
)

func main() {
	logging.InitLogging()
	defer logging.Logger.Sync()

	database.InitDBClient()
	cmd.Execute()

	r := routes.SetupRoutes()
	r.Run(":3000")
}
