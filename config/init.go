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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"io/ioutil"
	"metatds/utils"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/predatorpc/durafmt"
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

	// issue with too many open files
	http.DefaultClient.Timeout = time.Second * 10

	// цепляем редис и потом, проверяем постоянно, как у него дела
	RedisDBChan()

	// Напишем всем, что мы стартанули
	tlgrmRecipients := utils.Explode(Cfg.Telegram.Recipients, "; ")
	tlgrm := Telegram.Init(tlgrmRecipients, Cfg.Telegram.Socks5User, Cfg.Telegram.Socks5Password,
		Cfg.Telegram.Socks5Proxy, Cfg.Telegram.ApiURL, Cfg.Telegram.Token, Cfg.Telegram.UseProxy)

	timeStamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		UpTime.Year(), UpTime.Month(), UpTime.Day(), UpTime.Hour(), UpTime.Minute(), UpTime.Second())

	Telegram.SendMessage("\n" + timeStamp + "\n" + Cfg.General.Name + "\nTDS Service started\n")

	if tlgrm {
		if Cfg.Debug.Level > 0 {
			utils.PrintInfo("Telegram", "Successfully init Telegram Adapter", initModuleName)
		}
	} else {
		utils.PrintError("Error", "Init Telegram Adapter", initModuleName)
	}

	// начинаем слать статистику
	TDSStatisticChan()

	// начинаем перезагружать конфиг
	ReloadConfigChan()

	// отправка кликов в мета-дату
	RedisSendOrSaveClicks()

	// File sender это ресенд если не удалось предыдущее
	SendFileToRecieveApi()
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
			// check connection via Pong
			pong, err := Redisdb.Ping().Result()

			if err != nil {
				IsRedisAlive = false
				msg := "Can't connect to Redis server at " + Cfg.Redis.Host + ":" + strconv.Itoa(Cfg.Redis.Port)

				if Cfg.Debug.Level > 0 {
					utils.PrintError("Redis error", msg, initModuleName)
				}

				time.Sleep(60 * time.Second) // поспим чуть чуть

				goto tryUntilConnect
				//		os.Exit(0)
			} else {
				if Cfg.Debug.Level > 0 && !IsRedisAlive {
					utils.PrintSuccess("Redis response", pong, initModuleName)
				}

				if Cfg.Debug.Level > 0 && !IsRedisAlive {
					utils.PrintDebug("Completed", err, initModuleName)
				}

				IsRedisAlive = true
			}

			// defer runtime.GC()
			time.Sleep(60 * time.Second) // поспим чуть чуть
		}
	}()

	return c
}

func RedisSendOrSaveClicks() <-chan string {
	c := make(chan string)

	go func() {
		for {
			if IsRedisAlive {
				var clicks []map[string]string

				t := time.Now()
				timestamp := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
					t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

				timestampPrintable := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
					t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

				keys, _ := Redisdb.Keys("*:click:*").Result()

				for _, item := range keys {
					d, _ := Redisdb.HGetAll(item).Result()
					clicks = append(clicks, d)
				}

				TDSStatistic.ClicksSentToRedis += len(clicks)

				jsonData, _ := json.Marshal(clicks)

				if Cfg.Debug.Level > 1 {
					fmt.Println("Time elapsed export: ", time.Since(t))
				}

				if len(jsonData) > 0 && len(clicks) > 0 {

					url := Cfg.Click.ApiUrl     // "http://116.202.27.130/set/hits"
					token := Cfg.Click.ApiToken // "PaILgFTQQCvX9tzS"
					req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
					req.Header.Set("X-Token", token)
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Connection", "close")

					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						recover()
						goto tryagain
					}
					defer resp.Body.Close()

					utils.PrintInfo("Response status", resp.Status, initModuleName)

					if resp.Status == "200 OK" {
						body, _ := ioutil.ReadAll(resp.Body)
						utils.PrintInfo("Response", string(body), initModuleName)

						for _, item := range keys {
							_ = Redisdb.Del(item).Err()
						}

						if Cfg.Debug.Level > 1 {
							Telegram.SendMessage("\n" + timestampPrintable + "\n" +
								Cfg.General.Name + "\nClicks sent to API: " + strconv.Itoa(TDSStatistic.ClicksSentToRedis) +
								"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))
						}
					} else {
						utils.CreateDirIfNotExist("clicks")
						ioutil.WriteFile("clicks/"+timestamp+".json", jsonData, 0777)

						for _, item := range keys {
							_ = Redisdb.Del(item).Err()
						}

						Telegram.SendMessage("\n" + timestampPrintable + "\n" +
							Cfg.General.Name + "\nClicks saved to file" +
							"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))
					}
				}

				if Cfg.Debug.Level > 1 {
					fmt.Println("Time elapsed total: ", time.Since(t))
				}

				//defer runtime.GC()

			}
		tryagain:
			time.Sleep(time.Duration(1+Cfg.Click.DropToRedis) * time.Minute)
		}
	}()
	return c
}

func GetSystemConfiguration() string {
	text := spew.Sdump(Cfg)
	return text
}

//
// Send -*.json stored in files to reciever API
//

func SendFileToRecieveApi() <-chan string {
	c := make(chan string)
	go func() {
		for {
			var fdsReplace string

			t := time.Now()
			timestampPrintable := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
				t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

			fds, _ := filepath.Glob("clicks/*.json")

			if len(fds) > 0 {

				for _, item := range fds {

					fdsReplace = filepath.Base(item)

					// возможно надо проверку, но не уверен
					fileData, _ := ioutil.ReadFile(item)

					url := Cfg.Click.ApiUrl     // "http://116.202.27.130/set/hits"
					token := Cfg.Click.ApiToken // "PaILgFTQQCvX9tzS"
					req, err := http.NewRequest("POST", url, bytes.NewBuffer(fileData))
					req.Header.Set("X-Token", token)
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Connection", "close")

					client := &http.Client{}
					resp, err := client.Do(req)

					if resp != nil {
						recover()
						// TODO this needs to be recovered from panic otherwise fails
						fmt.Fprintln(os.Stderr, "can't GET page:", err)
					}

					defer resp.Body.Close()

					utils.PrintDebug("Response status", resp.Status, initModuleName)

					if resp.Status == "200 OK" {
						body, _ := ioutil.ReadAll(resp.Body)
						utils.PrintInfo("Response", string(body), initModuleName)

						// удаляем файл, мы его успешно обработали
						os.Remove(item)

						Telegram.SendMessage("\n" + timestampPrintable + "\n" +
							Cfg.General.Name + "\nResending file succedeed " + fdsReplace + " to API" +
							"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))
					} else {
						utils.PrintDebug("Error", "Sending file to click API failed", initModuleName)

						Telegram.SendMessage("\n" + timestampPrintable + "\n" +
							Cfg.General.Name + "\nResending file failed " + fdsReplace + " to API" +
							"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))
					}

					// поспим между файлами
					time.Sleep(time.Second * 10)
				}
			}

			// defer runtime.GC()
			time.Sleep(time.Duration(Cfg.Click.DropFilesToAPI) * time.Minute)
		}
	}()
	return c
}

func GetSystemStatistics() string {
	var text = "no stat"

	if TDSStatistic != (utils.TDSStats{}) {
		t := time.Now()
		timeStamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

		var memory runtime.MemStats
		var duration time.Duration // current duration & uptime
		var uptime, processingTime, memoryUsageGeneral, memoryUsagePrivate, avgReq string
		var openedFiles = "0"

		duration = 60 * time.Minute

		//if TDSStatistic.ProcessingTime < duration {
		if time.Since(UpTime) < duration {
			uptime = durafmt.Parse(time.Since(UpTime)).String(durafmt.DF_LONG)
			processingTime = durafmt.Parse(TDSStatistic.ProcessingTime).String(durafmt.DF_LONG)
		} else {
			uptime = durafmt.Parse(time.Since(UpTime)).String(durafmt.DF_MIDDLE)
			processingTime = durafmt.Parse(TDSStatistic.ProcessingTime).String(durafmt.DF_MIDDLE)
		}

		runtime.ReadMemStats(&memory)

		RealDetectedGeneral := memory.Sys + memory.HeapSys + memory.HeapAlloc + memory.HeapInuse - memory.Alloc
		RealDetectedPrivate := memory.HeapSys - memory.Alloc

		memoryUsageGeneral = strconv.FormatUint(utils.BToMb(RealDetectedGeneral), 10)
		memoryUsagePrivate = strconv.FormatUint(utils.BToMb(RealDetectedPrivate), 10)

		//fmt.Print("[MEMORY USAGE]",memoryUsage, memory.Sys)

		if Cfg.General.OS == "linux" || Cfg.General.OS == "unix" {

			pid := strconv.Itoa(os.Getpid())
			fds, e := ioutil.ReadDir("/proc/" + pid + "/fd")

			if e != nil && Cfg.Debug.Level > 0 {
				utils.PrintError("Error", "reading process directory failed", initModuleName)
			} else {
				//utils.PrintInfo("PID", pid, initModuleName)
			}

			if len(fds) > 0 {
				openedFiles = strconv.Itoa(len(fds))
			}
		}

		dur:=DurationAverage(utils.ResponseAverage)

		if dur < time.Duration(1 * time.Millisecond) {
			avgReq = durafmt.Parse(dur).String(durafmt.DF_LONG)
		} else {
			avgReq = " < 1 ms"
		}

		uniqueRequests := TDSStatistic.RedirectRequest - TDSStatistic.CookieRequest - TDSStatistic.IncorrectRequest

		text = "\n" + timeStamp + "\n" + Cfg.General.Name +
			"\n\nINFO" +
			"\n\nFlow update request    : " + strconv.Itoa(TDSStatistic.UpdatedFlows) +
			"\nFlow appended          : " + strconv.Itoa(TDSStatistic.AppendedFlows) +
			//"\nPixel request          : " + strconv.Itoa(TDSStatistic.PixelRequest) +
			"\nClick Info request     : " + strconv.Itoa(TDSStatistic.ClickInfoRequest) +
			"\nFlow Info request      : " + strconv.Itoa(TDSStatistic.FlowInfoRequest) +
			"\nRedirect request       : " + strconv.Itoa(TDSStatistic.RedirectRequest) +
			//			"\nRedis Stat request     : " + strconv.Itoa(TDSStatistic.RedisStatRequest) +
			"\nIncorrect request      : " + strconv.Itoa(TDSStatistic.IncorrectRequest) +
			"\nCookies request        : " + strconv.Itoa(TDSStatistic.CookieRequest) +
			"\nUnique request (?)     : " + strconv.Itoa(uniqueRequests) +
			"\n\nUp time                : " + uptime +
			"\nProcessing time        : " + processingTime +
			"\nAverage response time  : " + avgReq +
			"\n\nSYSTEM INFO" +
			"\n\nOperating system       : " + Cfg.General.OS +
			"\nDebug level            : " + strconv.Itoa(Cfg.Debug.Level) +
			"\nTotal memory allocated : " + memoryUsageGeneral + " Mb" +
			"\nPrivate memory         : " + memoryUsagePrivate + " Mb" +
			"\nOpened files           : " + openedFiles +
			"\n\nREDIS" +
			"\n\nConnection             : " + strconv.FormatBool(IsRedisAlive) +
			"\nClicks sent            : " + strconv.Itoa(TDSStatistic.ClicksSentToRedis) +
			"\n"
		return text
	} else {
		TDSStatistic.Reset()
		return text
	}
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
				text := GetSystemStatistics()
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

			//defer runtime.GC()
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
//
//
// func DurationAverage(dur []time.Duration) time.Duration {
// 	var allTime time.Duration
// 	for _, item := range dur {
// 		allTime += item
// 	}
// 	//division by zero
// 	fmt.Println("All time sum ", allTime, " / ", len(dur))
//
// 	return allTime / 1 + time.Duration(len(dur)) //time.Duration(1+len(dur))
// }

func DurationAverage(dur []time.Duration) time.Duration {
	var allTime float64
	for _, item := range dur {
		allTime += float64(item)
	}
	result:= allTime / float64(1+len(dur))
//	fmt.Println("All time sum ", allTime, " / ", 1+len(dur))
	return time.Duration(result)
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
			//defer runtime.GC()
			time.Sleep(time.Duration(1+Cfg.General.ConfReload*60) * time.Second)
		}
	}()

	return c
}
