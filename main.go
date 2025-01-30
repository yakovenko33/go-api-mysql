package main

import (
	"go-api-docker/database"
	"go-api-docker/routes"
)

// goose  create new_user_table sql
// goose -dir ./db/migrations mysql "user:password@tcp(mysql:3306)/my_database?parseTime=true" up
// goose mysql "user:password@tcp(mysql:3306)/my_database?parseTime=true" up / down
func main() {
	database.InitDBClient()
	r := routes.SetupRoutes()
	r.Run(":3000")
}
