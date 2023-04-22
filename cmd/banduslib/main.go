package main

import (
	"banduslib/internal/api"
	"banduslib/internal/database"
	"banduslib/internal/utils"
	"fmt"
	"github.com/rs/zerolog/log"
)

func main() {
	utils.SetupConsoleLogger()

	db, err := database.GetInitSqliteDb("main.db")
	if err != nil {
		panic(fmt.Sprintf("failed initializing database: %s", err.Error()))
	}

	dbService := database.NewDbDataService(db)
	apiHandler := api.NewApiHandler(dbService)
	apiRouter := api.NewApiRouter(apiHandler)

	router := apiRouter.GetEngine()

	log.Fatal().Err(router.Run(":8081"))
}
