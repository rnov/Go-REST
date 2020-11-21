package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rnov/Go-REST/pkg/auth"
	infra "github.com/rnov/Go-REST/pkg/config"
	"github.com/rnov/Go-REST/pkg/db"
	"github.com/rnov/Go-REST/pkg/http/rest"
	"github.com/rnov/Go-REST/pkg/logger"
	"github.com/rnov/Go-REST/pkg/service"
)

const (
	EnvVarPath = "ENV_PATH"
)

func main() {
	fmt.Println("Hello, 世界")

	// load the application configuration
	envConfigPath, present := os.LookupEnv(EnvVarPath)
	if !present {
		log.Fatal("Env Variable Not present")
	} else if len(envConfigPath) == 0 {
		log.Fatal("Empty Env Variable")
	}

	cfg, err := infra.LoadAPIConfig(envConfigPath)

	if err != nil {
		log.Fatal(err)
	}

	// create RecipeSrv custom logger
	//l := logger.NewLogger(cfg.RedisLog, cfg.LogsPath)
	l := logger.NewLogger()
	//fmt.Sprint(l)

	// create DB client
	dbClient, err := db.NewClient(cfg.DBCfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	// get auth accessor
	authorization := auth.NewAuth(dbClient, l)

	// initialize controllers
	// In this case recipe and rate share same DB and logger but could be different ones
	RecipeSrv := service.NewRecipe(dbClient)
	RateSrv := service.NewRate(dbClient)

	// Create handlers
	rcpHandler := rest.NewRecipeHandler(RecipeSrv, l)
	rateHandler := rest.NewRateHandler(RateSrv, l)

	r := rest.NewRouter(rcpHandler, rateHandler, authorization)

	// Fire up the server
	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))
}
