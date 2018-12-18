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
	"github.com/go-redis/redis"
	"metatds/utils"
	"net/url"
	"strconv"
	"time"

	"github.com/hako/durafmt"
)

const initModuleName = "init.go"

var UpTime time.Time

func init() {
	// get current timestamp
	UpTime = time.Now()

	// сначала загружаем настройки потом, цепляем все остальное
	InitConfig()

	if Cfg.Debug.Level > 0 {
		utils.PrintDebug("Initialization", "", initModuleName)
	}

	// цепляем редис и потом, проверяем постоянно, как у него дела
	RedisDBChan()

	// Напишем всем, что мы стартанули
	tlgrmRecipients := utils.Explode(Cfg.Telegram.Recipients, "; ")
	tlgrm := Telegram.Init(tlgrmRecipients, Cfg.Telegram.Socks5User, Cfg.Telegram.Socks5Password,
		Cfg.Telegram.Socks5Proxy, Cfg.Telegram.ApiURL, Cfg.Telegram.Token)

	timeStamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		UpTime.Year(), UpTime.Month(), UpTime.Day(), UpTime.Hour(), UpTime.Minute(), UpTime.Second())

	Telegram.SendMessage(Cfg.General.Name + ": TDS Service started @ " + timeStamp)

	if tlgrm {
		if Cfg.Debug.Level > 0 {
			utils.PrintInfo("Telegram", "Successfully init Telegram Adapter", initModuleName)
		}
	} else {
		utils.PrintError("Error", "Init Telegram Adapter", initModuleName)
	}

	// начинаем считать статистику
	TDSStatisticChan()

	// начинаем перезагружать конфиг
	ReloadConfigChan()

}

/*
*
* Init Redis connection /Reload Redis connection if broken
*
 */

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
					utils.PrintSuccess("Config", msg, initModuleName)
				}
			}

		tryUntilConnect: // try to reconnect until success
			// start responding json {error:redis}

			// check connection via Pong
			pong, err := Redisdb.Ping().Result()

			if err != nil {
				IsRedisAlive = false
				msg := "Can't connect to Redis server at " + Cfg.Redis.Host + ":" + strconv.Itoa(Cfg.Redis.Port)

				if Cfg.Debug.Level > 0 {
					utils.PrintError("Redis error", msg, initModuleName)
				}

				goto tryUntilConnect
				//		os.Exit(0)
			} else {
				if Cfg.Debug.Level > 0 && !IsRedisAlive {
					utils.PrintSuccess("Redis response", pong, initModuleName)
					//utils.PrintSuccess("Redis response", err, initModuleName)
				}

				if Cfg.Debug.Level > 0 && !IsRedisAlive {
					utils.PrintDebug("Completed", err, initModuleName)
				}

				IsRedisAlive = true
			}

			time.Sleep(10 * time.Second) // поспим чуть чуть
		}
	}()

	return c
}

/*
*
* Reload config channel
*
 */

func ReloadConfigChan() <-chan string {
	c := make(chan string)

	go func() {
		for {
			// перезагружаем конфиг и идем спать
			ReloadConfig()
			// поспим чуть чуть
			// +1 its to avoid dumbs with zero multiplication
			time.Sleep(time.Duration(1+Cfg.General.ConfReload*60) * time.Second)
		}
	}()

	return c
}

/*
*
* Telegram send statistic channel
*
 */

func TDSStatisticChan() <-chan string {
	c := make(chan string)

	go func() {
		for {

			if TDSStatistic != (utils.TDSStats{}) {
				t := time.Now()
				timeStamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
					t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

				text := Cfg.General.Name + " usage\n" + timeStamp + ":" +
					"\n\nUpdate Flow: " + strconv.Itoa(TDSStatistic.UpdatedFlows) +
					"\nAppende Flow: " + strconv.Itoa(TDSStatistic.AppendedFlows) +
					"\nPixel Request: " + strconv.Itoa(TDSStatistic.PixelRequest) +
					"\nClick Info Request: " + strconv.Itoa(TDSStatistic.ClickInfoRequest) +
					"\nFlow Info Request: " + strconv.Itoa(TDSStatistic.FlowInfoRequest) +
					"\nRedirect Request: " + strconv.Itoa(TDSStatistic.RedirectRequest) +
					"\nRedis Stat Request: " + strconv.Itoa(TDSStatistic.RedisStatRequest) +
					"\nIncorrect Request: " + strconv.Itoa(TDSStatistic.IncorrectRequest) +
					"\n\nUp time: " + durafmt.Parse(time.Since(UpTime)).String() +
					"\nWorkt time: " + durafmt.Parse(TDSStatistic.WorkTime).String() +
					"\n\nRedis connection: " + strconv.FormatBool(IsRedisAlive)

				if Telegram.SendMessage(url.QueryEscape(text)) {
					if Cfg.Debug.Level > 0 {
						utils.PrintInfo("Telegram", "Sending message success", initModuleName)
					}
				} else {
					if Cfg.Debug.Level > 0 {
						utils.PrintError("Telegram", "Sending message error", initModuleName)
					}
				}
			} else {
				TDSStatistic.Reset()
			}

			// +1 its to avoid dumbs with zero multiplication
			time.Sleep(time.Duration(1+Cfg.Telegram.MsgInterval*60) * time.Second) // поспим чуть чуть
		}
	}()

	return c
}
