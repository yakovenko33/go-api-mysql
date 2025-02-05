package main

import (
	"go-api-docker/GoCrm/Common/Logging"
	"go-api-docker/cmd"
	"go-api-docker/database"
	"go-api-docker/routes"
)

func main() {
	Logging.InitLogging()
	database.InitDBClient()
	cmd.Execute()
	r := routes.SetupRoutes()
	r.Run(":3000")
}
