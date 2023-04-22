package main

import (
	"banduslib/internal/api"
	"banduslib/internal/database"
	"fmt"
	"log"
)

func main() {
	db, err := database.GetInitSqliteDb("main.db")
	if err != nil {
		panic(fmt.Sprintf("failed initializing database: %s", err.Error()))
	}

	dbService := database.NewDbDataService(db)
	apiHandler := api.NewApiHandler(dbService)
	apiRouter := api.NewApiRouter(apiHandler)

	router := apiRouter.GetEngine()

	log.Fatal(router.Run(":8081"))
}
