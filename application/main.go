package main

import (
	"Go-REST/application/apiRest"
	"Go-REST/application/apiRest/handlers"
	infra "Go-REST/application/common"
	"Go-REST/application/controller"
	"Go-REST/application/dbInterface"
	"fmt"
	"log"
	"net/http"
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

	// create RecipeController custom logger
	logger := infra.NewLogger(cfg.RedisLog, cfg.LogsPath)

	// get redis master (and only one)
	redisMaster := cfg.Redis[0]
	redisClient := dbInterface.NewClient(redisMaster.Host, redisMaster.Port, redisMaster.Db)

	// check connection with redis
	pong, err := redisClient.Ping().Result()
	fmt.Println(pong)
	if err != nil {
		panic(err)
	}

	// create redisProxy with the given client (master)
	redisProxy := dbInterface.NewRedisProxy(redisClient)

	// initialize controllers
	// In this case recipe and rate share same DB and logger but could be different ones
	RecipeController := controller.NewRecipeController(redisProxy, logger)
	RateController := controller.NewRateController(redisProxy, logger)

	// Create handlers
	rcpHandler := handlers.NewRecipeHandler(RecipeController)
	rateHandler := handlers.NewRateHandler(RateController)

	r := apiRest.NewRouter(rcpHandler, rateHandler)

	// Fire up the server
	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))

}
