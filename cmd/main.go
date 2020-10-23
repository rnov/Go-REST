package main

import (
	"github.com/rnov/Go-REST/pkg/db"
	"github.com/rnov/Go-REST/pkg/http/rest"
	"github.com/rnov/Go-REST/pkg/rate"
	"github.com/rnov/Go-REST/pkg/recipe"
	"net/http"

	"fmt"
	infra "github.com/rnov/Go-REST/pkg/config"
	"github.com/rnov/Go-REST/pkg/logger"
	"log"
	"os"
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

	cfg, err := infra.LoadApiConfig(envConfigPath)

	if err != nil {
		log.Fatal(err)
	}

	// create RecipeSrv custom logger
	l := logger.NewLogger(cfg.RedisLog, cfg.LogsPath)
	fmt.Sprint(l)

	// fixme actually the param should be config struct with some data like host, port etc ...
	dbClient, err := db.NewDbClient("redis")
	if err != nil {
		log.Fatal(err.Error())
	}

	// initialize controllers
	// In this case recipe and rate share same DB and logger but could be different ones
	RecipeSrv := recipe.NewRecipeSrv(dbClient, l)
	RateSrv := rate.NewRateSrv(dbClient, l)

	// Create handlers
	rcpHandler := rest.NewRecipeHandler(RecipeSrv)
	rateHandler := rest.NewRateHandler(RateSrv)

	r := rest.NewRouter(rcpHandler, rateHandler)

	// Fire up the server
	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))

}
