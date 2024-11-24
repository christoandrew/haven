package server

import (
	"github.com/christo-andrew/haven/pkg/config"
	"github.com/christo-andrew/haven/pkg/database"
	"log"

	"github.com/christo-andrew/haven/internal/api"
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
	currentConfig, err := config.New(".env")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	currentConfig.Validate()
	app := api.NewApiServer(currentConfig)
	db, err := currentConfig.Database.GetDB()

	if err != nil {
		log.Fatal(err)
	}

	server := app.SetupRouter(db)
	database.Migrate(db)
	err = server.Run()
	if err != nil {
		panic(err)
	}
}
