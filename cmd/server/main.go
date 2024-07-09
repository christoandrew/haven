package server

import (
	"log"

	"github.com/christo-andrew/haven/internal/api"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/config"
	database "github.com/christo-andrew/haven/pkg/database"
	"github.com/joho/godotenv"
)

//	@title			Haven API
//	@version		1.0

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	andrew.christo.wekesa@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

func StartApp() {
	serverConfig := config.DefaultServerConfig()
	err := godotenv.Load(serverConfig.EnvPath)
	if err != nil {
		log.Fatal(err)
	}
	db := database.GetDB(serverConfig.DatabaseConfig)
	app := api.NewServer(serverConfig)
	server := app.SetupRouter(db)
	models.Migrate(db)
	err = server.Run()
	if err != nil {
		panic(err)
	}
}
