package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

const initModuleName = "init.go"

func init() {

	initConfig() // loading configuration globally

	if cfg.Debug.Level > 0 {
		printDebug("Initialization", "", initModuleName)
	}

	// get connection to Redis
	redisdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password, // password set
		DB:       0,                  // use default DB
	})

	//verifying config
	if cfg.Debug.Level > 0 {
		msg := " Redis = " + cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port) + ", Self = " +
			cfg.General.Host + ":" + strconv.Itoa(cfg.General.Port)

		printSuccess("Config", msg, initModuleName)
	}

	// check connection via Pong
	pong, err := redisdb.Ping().Result()

	if err != nil {
		msg := "Can't connect to Redis server at " + cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port)
		printError("Redis error", msg, initModuleName)
		os.Exit(0)

	} else {
		if cfg.Debug.Level > 0 {
			printSuccess("Redis response", pong, initModuleName)
			printSuccess("Redis response", err, initModuleName)
		}

		if cfg.Debug.Level > 0 {
			printDebug("Completed", err, initModuleName)
		}
	}
}
