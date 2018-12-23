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
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"io/ioutil"
	"metatds/utils"
	"runtime"
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

	Telegram.SendMessage("```\n" + timeStamp + "\n" + Cfg.General.Name + "\nTDS Service started\n```")

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

	// this is temporary workaround, should be removed on release
	// TODO remove in release
	//TempResetRedisClicks()
	RedisSaveClicks()

	// TODO channel to resend stats to API, when Cherkesov gives me an URL
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

func RedisSaveClicks() <-chan string {
	c := make(chan string)

	go func() {
		for {
			var clicks []map[string]string

			t := time.Now()
			timestamp := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
				t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

			keys, _ := Redisdb.Keys("*:click:*").Result()
			// 	fmt.Println("Keys found by mask: ", len(keys))

			for _, item := range keys {
				d, _ := Redisdb.HGetAll(item).Result()
				clicks = append(clicks, d)
			}

			jsonData, _ := json.Marshal(clicks)

			if Cfg.Debug.Level > 1 {
				fmt.Println("Time elapsed export: ", time.Since(t))
			}

			if len(jsonData) > 0 && len(clicks) > 0 {
				utils.CreateDirIfNotExist("clicks")
				ioutil.WriteFile("clicks/"+timestamp+".json", jsonData, 0777)

				// TODO обработка если в файл не записалось пока просто грохаем
				for _, item := range keys {
					_ = Redisdb.Del(item).Err()
				}
			}
			timestampPrintable := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
				t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

			if Cfg.Debug.Level > 1 {
				fmt.Println("Time elapsed total: ", time.Since(t))
			}

			Telegram.SendMessage("```\n" + timestampPrintable + "\n" +
				Cfg.General.Name + "\nClicks saved and reseted from RedisDB\n```")
			time.Sleep(time.Duration(1+Cfg.Click.DropToRedis) * time.Minute)
		}
	}()
	return c
}

func TempResetRedisClicks() <-chan string {
	c := make(chan string)

	go func() {
		for {
			start := time.Now()

			keys, _ := Redisdb.Keys("*:click:*").Result()
			fmt.Println("Keys found by mask: ", len(keys))

			for _, item := range keys {
				_ = Redisdb.Del(item).Err()
			}

			fmt.Println("Time elapsed: ", time.Since(start))

			t := time.Now()
			timeStamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
				t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

			Telegram.SendMessage("```\n" + timeStamp + "\n" + Cfg.General.Name + "\nRedis reset clicks succeeded\n```")

			time.Sleep(1 * time.Hour)
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

				var memory runtime.MemStats
				var duration time.Duration // current duration & uptime
				var uptime, processingTime, memoryUsage string

				duration = 10 * time.Minute

				//if TDSStatistic.ProcessingTime < duration {
				if time.Since(UpTime) < duration {
					uptime = durafmt.Parse(time.Since(UpTime)).String(durafmt.DF_LONG)
					processingTime = durafmt.Parse(TDSStatistic.ProcessingTime).String(durafmt.DF_LONG)
				} else {
					uptime = durafmt.Parse(time.Since(UpTime)).String(durafmt.DF_SHORT)
					processingTime = durafmt.Parse(TDSStatistic.ProcessingTime).String(durafmt.DF_SHORT)
				}

				runtime.ReadMemStats(&memory)
				memoryUsage = strconv.FormatUint(utils.BToMb(memory.Sys), 10)
				//fmt.Print("[MEMORY USAGE]",memoryUsage, memory.Sys)

				uniqueRequests := TDSStatistic.RedirectRequest - TDSStatistic.CookieRequest - TDSStatistic.IncorrectRequest

				text := "```\n" + timeStamp + "\n" + Cfg.General.Name + " usage:" +
					"\n\nUpdate flow       : " + strconv.Itoa(TDSStatistic.UpdatedFlows) +
					"\nAppende flow      : " + strconv.Itoa(TDSStatistic.AppendedFlows) +
					"\nPixel request     : " + strconv.Itoa(TDSStatistic.PixelRequest) +
					"\nClick Info request: " + strconv.Itoa(TDSStatistic.ClickInfoRequest) +
					"\nFlow Info request : " + strconv.Itoa(TDSStatistic.FlowInfoRequest) +
					"\nRedirect request  : " + strconv.Itoa(TDSStatistic.RedirectRequest) +
					"\nRedis Stat request: " + strconv.Itoa(TDSStatistic.RedisStatRequest) +
					"\nIncorrect request : " + strconv.Itoa(TDSStatistic.IncorrectRequest) +
					"\nCookies request   : " + strconv.Itoa(TDSStatistic.CookieRequest) +
					"\nUnique request    : " + strconv.Itoa(uniqueRequests) +
					"\n\nUp time           : " + uptime +
					"\nProcessing time   : " + processingTime +
					"\nMemory allocated  : " + memoryUsage + " Mb" +
					"\n\nRedis connection  : " + strconv.FormatBool(IsRedisAlive) + "\n```"

				if Telegram.SendMessage(text) {
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

//
// Template for Channel by predator_pc
//
// func ChannelWithSleepTemplate() <-chan string {
// 	c := make(chan string)
// 	go func() {
// 		for {
// 			time.Sleep(time.Minute * 10)
// 		}
// 	}()
// 	return c
// }
