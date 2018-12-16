/****************************************************************************************************
*
* Initialization package preparing to do hard work :)) special for Meta CPA, Ltd.
* by Michael S. Merzlyakov AFKA predator_pc@12122018
* version v2.0.3
*
* created at 04122018
* last edit: 16122018
*
*****************************************************************************************************/

package config

import (
	"fmt"
	"metatds/utils"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const initModuleName = "init.go"

func init() {

	InitConfig() // loading configuration globally

	tlgrmRecipients := utils.Explode(Cfg.Telegram.Recipients, "; ")
	tlgrm := Telegram.Init(tlgrmRecipients, Cfg.Telegram.Socks5User, Cfg.Telegram.Socks5Password,
		Cfg.Telegram.Socks5Proxy, Cfg.Telegram.ApiURL, Cfg.Telegram.Token)

	t := time.Now()
	timeStamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	Telegram.SendMessage(Cfg.General.Name + ": TDS Service started @ " + timeStamp)

	if tlgrm {
		if Cfg.Debug.Level > 0 {
			utils.PrintInfo("Telegram", "Successfully init Telegram Adapter", initModuleName)
		}
	} else {
		utils.PrintError("Error", "Init Telegram Adapter", initModuleName)
	}

	if Cfg.Debug.Level > 0 {
		utils.PrintDebug("Initialization", "", initModuleName)
	}

	// get connection to Redis
	Redisdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Cfg.Redis.Host, Cfg.Redis.Port),
		Password: Cfg.Redis.Password, // password set
		DB:       0,                  // use default DB
	})

	//verifying config
	if Cfg.Debug.Level > 0 {
		msg := " Redis = " + Cfg.Redis.Host + ":" + strconv.Itoa(Cfg.Redis.Port) + ", Self = " +
			Cfg.General.Host + ":" + strconv.Itoa(Cfg.General.Port)

		utils.PrintSuccess("Config", msg, initModuleName)
	}

	// check connection via Pong
	pong, err := Redisdb.Ping().Result()

	if err != nil {
		msg := "Can't connect to Redis server at " + Cfg.Redis.Host + ":" + strconv.Itoa(Cfg.Redis.Port)
		utils.PrintError("Redis error", msg, initModuleName)
		os.Exit(0)

	} else {
		if Cfg.Debug.Level > 0 {
			utils.PrintSuccess("Redis response", pong, initModuleName)
			utils.PrintSuccess("Redis response", err, initModuleName)
		}

		if Cfg.Debug.Level > 0 {
			utils.PrintDebug("Completed", err, initModuleName)
		}
	}
}
