package main

import (
	"fmt"
	"os"

	"go.uber.org/fx"

	cli "go-api-docker/cmd/cli"
	database "go-api-docker/internal/common/database"
	logging "go-api-docker/internal/common/logging"
	access_control_model "go-api-docker/internal/common/security/access_control_models"
	server "go-api-docker/internal/common/server"
	register_routes "go-api-docker/internal/common/server/controllers"
)

func coreApp() *fx.App {
	app := fx.New(
		fx.Provide(
			logging.InitLogging,
			database.ProvideDBConnection,
			access_control_model.InitAccessControlModel,
			server.NewRouter,
		),
		register_routes.Controllers,
		fx.Invoke(
			server.StartServer,
		),
	)
	return app
}

func main() {
	if len(os.Args) < 2 {
		app := coreApp()
		app.Run()
		return
	}

	switch os.Args[1] {
	case "cli":
		os.Args = os.Args[1:]
		cli.RunCLI()
	default:
		fmt.Println("Unknown mode. Usage: myapp <web|cli>")
	}
}
