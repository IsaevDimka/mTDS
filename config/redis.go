package config

import (
	"fmt"
	"github.com/go-redis/redis"
	"metatds/utils"
	"strconv"
	"time"
)

const redisModuleName = "redis.go"

//
// Init Redis connection /Reload Redis connection if broken
//
func RedisDBChan() <-chan string {
	c := make(chan string)

	go func() {
		// get connection to Redis
		Redisdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", Cfg.Redis.Host, Cfg.Redis.Port),
			Password: Cfg.Redis.Password, // password set
			DB:       0,                  // use default DB
		})
		defer Redisdb.Close() // Если редис отвалится, потом конекция не повиснет

		for {

			// verifying config
			// напечатаем его заодним на экран
			if Cfg.Debug.Level > 0 && !IsRedisAlive {
				msg := " Redis = " + Cfg.Redis.Host + ":" + strconv.Itoa(Cfg.Redis.Port) + ", Self = " +
					Cfg.General.Host + ":" + strconv.Itoa(Cfg.General.Port)
				if Cfg.Debug.Level > 0 && !IsRedisAlive {
					utils.PrintSuccess("Config", msg, redisModuleName)
				}
			}

		tryUntilConnect: // try to reconnect until success
			// check connection via Pong

			pong, err := Redisdb.Ping().Result()

			if err != nil {
				IsRedisAlive = false
				msg := "Can't connect to Redis server at " + Cfg.Redis.Host + ":" + strconv.Itoa(Cfg.Redis.Port)

				if Cfg.Debug.Level > 0 {
					utils.PrintError("Redis error", msg, redisModuleName)
				}

				time.Sleep(10 * time.Second) // поспим чуть чуть

				goto tryUntilConnect
			} else {
				if Cfg.Debug.Level > 0 && !IsRedisAlive {
					utils.PrintSuccess("Redis response", pong, redisModuleName)
				}

				if Cfg.Debug.Level > 0 && !IsRedisAlive {
					utils.PrintDebug("Completed", err, redisModuleName)
				}

				IsRedisAlive = true
			}

			// defer runtime.GC()
			time.Sleep(10 * time.Second) // поспим чуть чуть
		}
	}()

	return c
}
