package main

import (
	"go-api-docker/cmd"
	"go-api-docker/database"
	"go-api-docker/routes"
)

func main() {
	database.InitDBClient()
	cmd.Execute()
	r := routes.SetupRoutes()
	r.Run(":3000")
}
