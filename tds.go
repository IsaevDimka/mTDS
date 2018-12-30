/****************************************************************************************************
*
* Main TDS module, special for Meta CPA, Ltd.
* by Michael S. Merzlyakov AFKA predator_pc@02122018
* version v2.0.3
*
* created at 04122018
* last edit: 16122018
*
* usage: $ tds run
*
*****************************************************************************************************/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/predatorpc/durafmt"
	"github.com/sevenNt/echo-pprof"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"metatds/config"
	"metatds/models"
	"metatds/utils"
)

const tdsModuleName = "tds.go"
const timestampFile = "last.update.time"

// Карта ключей которые мы хотим получить
var keyMap = []string{"flow_hash", "click_hash", "sub1", "sub2", "sub3", "sub4", "sub5", "format",
	"click_id", "flow_id", "preland_id", "land_id"} // support for old version of TDS

var pixel = []byte(`data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HgAGgwJ/lK3Q6wAAAABJRU5ErkJggg==`)

// Тип для
type InfoData struct {
	Flow  models.FlowData  // модель потока
	Click models.ClickData // модель клика
}

/*
*
* Main GO package handler startup and settings
*
 */

func main() {

	//go utils.MemMonitor(1)
	UpdateFlowsListChan()

	// Echo instance
	router := echo.New()
	// Middleware
	//router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	//avoid chrome to request favicon
	router.GET("/favicon.ico", func(c echo.Context) error {
		return c.Blob(200, "image/png", pixel)
		//return c.String(404, "not found") //nothing
	})
	// Routes
	router.GET("/", flowHandler)
	router.GET("/:flow_hash", flowHandler)
	router.GET("/:flow_hash/", flowHandler)
	router.GET("/:flow_hash/:sub1", flowHandler)
	router.GET("/:flow_hash/:sub1/", flowHandler)
	router.GET("/:flow_hash/:sub1/:sub2", flowHandler)
	router.GET("/:flow_hash/:sub1/:sub2/", flowHandler)
	router.GET("/:flow_hash/:sub1/:sub2/:sub3", flowHandler)
	router.GET("/:flow_hash/:sub1/:sub2/:sub3/", flowHandler)
	router.GET("/:flow_hash/:sub1/:sub2/:sub3/:sub4", flowHandler)
	router.GET("/:flow_hash/:sub1/:sub2/:sub3/:sub4/", flowHandler)
	router.GET("/:flow_hash/:sub1/:sub2/:sub3/:sub4/:sub5", flowHandler)
	router.GET("/:flow_hash/:sub1/:sub2/:sub3/:sub4/:sub5/", flowHandler)
	// Routes
	router.GET("/r/", flowHandler)
	router.GET("/r/:flow_hash", flowHandler)
	router.GET("/r/:flow_hash/", flowHandler)
	router.GET("/r/:flow_hash/:sub1", flowHandler)
	router.GET("/r/:flow_hash/:sub1/", flowHandler)
	router.GET("/r/:flow_hash/:sub1/:sub2", flowHandler)
	router.GET("/r/:flow_hash/:sub1/:sub2/", flowHandler)
	router.GET("/r/:flow_hash/:sub1/:sub2/:sub3", flowHandler)
	router.GET("/r/:flow_hash/:sub1/:sub2/:sub3/", flowHandler)
	router.GET("/r/:flow_hash/:sub1/:sub2/:sub3/:sub4", flowHandler)
	router.GET("/r/:flow_hash/:sub1/:sub2/:sub3/:sub4/", flowHandler)
	router.GET("/r/:flow_hash/:sub1/:sub2/:sub3/:sub4/:sub5", flowHandler)
	router.GET("/r/:flow_hash/:sub1/:sub2/:sub3/:sub4/:sub5/", flowHandler)
	// Routes
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4/", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4/:sub5", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4/:sub5/", clickBuild)

	router.GET("/c/info/:click_hash", clickHandler)
	router.GET("/c/list", ListClickHandler)

	router.GET("/stat", GetSystemStatkHandler)
	router.GET("/conf", GetSystemConfHandler)

	customServer := &http.Server{
		Addr:         ":" + strconv.Itoa(config.Cfg.General.Port),
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
		IdleTimeout:  100 * time.Second,
	}

	customServer.SetKeepAlivesEnabled(false)

	router.HideBanner = true
	router.Logger.SetLevel(log.OFF)
	echopprof.Wrap(router)

	// run router
	if config.Cfg.General.Port != 0 {
		// regular server
		// router.Logger.Fatal(router.Start(":" + strconv.Itoa(config.Cfg.General.Port)))

		// custom server
		router.Logger.Fatal(router.StartServer(customServer))
	} else {
		//exit if not
		panic("[ERROR] Failed to obtain server port from settings.ini")
	}

}

func GetSystemConfHandler(c echo.Context) error {
	agent := c.Request().UserAgent()
	// защита от долбоебов
	if agent == "MetaDevAgent" {
		text := config.GetSystemConfiguration()
		return c.HTML(200, "<html><head><title>TDS System statistics</title><script>"+
			"setInterval(function(){window.location.reload(true)},5000);"+
			"</script></head><body><pre>"+
			text+"</pre></body></html>")

	} else {
		return c.String(404, "Not found on server")
	}
}

func GetSystemStatkHandler(c echo.Context) error {
	agent := c.Request().UserAgent()
	// защита от долбоебов
	if agent == "MetaDevAgent" {
		text := config.GetSystemStatistics()
		return c.HTML(200, "<html><head><title>TDS System statistics</title><script>"+
			"setInterval(function(){window.location.reload(true)},5000);"+
			"</script></head><body><pre>"+
			text+"</pre></body></html>")

	} else {
		return c.String(404, "Not found on server")
	}
}

func ImportFlowsToRedis(jsonData []byte) (int, bool) {
	Flows := make(map[string]models.FlowImportData)
	if err := json.Unmarshal(jsonData, &Flows); err != nil {
		if config.Cfg.Debug.Level > 1 {
			utils.PrintDebug("Error", "Can`t decode JSON given", tdsModuleName)
		}
	} else {
		for _, item := range Flows {
			_ = config.Redisdb.Set(item.Hash+":ID", item.ID, 0).Err()
			_ = config.Redisdb.Set(item.Hash+":Hash", item.Hash, 0).Err()
			_ = config.Redisdb.Set(item.Hash+":OfferID", item.OfferID, 0).Err()
			_ = config.Redisdb.Set(item.Hash+":WebMasterID", item.WebMasterID, 0).Err()
			_ = config.Redisdb.Set(item.Hash+":WebMasterCurrencyID", item.WebMasterCurrencyID, 0).Err()

			if len(item.Lands) > 0 {
				for i, lands := range item.Lands {
					_ = config.Redisdb.HSet(item.Hash+":land:"+strconv.Itoa(i), "id", lands.ID)
					_ = config.Redisdb.HSet(item.Hash+":land:"+strconv.Itoa(i), "url", lands.URL)
				}
			}
		}

		return len(Flows), true
	}
	return len(Flows), false
}

const defaultStartOfEpoch = "946684800"

func UpdateFlowsListChan() <-chan string {
	c := make(chan string)
	go func() {
		for {
			if config.IsRedisAlive {
				var body []byte

				t := time.Now()
				timestampWriteable := strconv.FormatInt(time.Now().Unix(), 10)
				timestampPrintable := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
					t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

				// getting current count of flows and if it isn't null then proceeed
				currentCount, _ := config.Redisdb.Keys("*:ID").Result()
				fileData, err := ioutil.ReadFile(timestampFile)

				if err == nil && len(currentCount) > 0 {
					// ---------------------------------------------------------------------------------------------------------------------
					// LOADING WITH PARAMS OF LAST UPDATE
					// When TDS starting up first time we need to load all flows in it
					// ---------------------------------------------------------------------------------------------------------------------
					// if we can't parse this we should get all anyway
					_, err := strconv.ParseInt(string(fileData), 10, 64)
					if err != nil && config.Cfg.Debug.Level > 0 {
						utils.PrintDebug("Error", "parsing timestamp from `"+timestampFile+"` failure", tdsModuleName)
						fileData = []byte(defaultStartOfEpoch) // 2000-01-01
					}

					// performing request to our API
					url := config.Cfg.Redis.ApiFlowsURL
					req, err := http.Get(url + strings.Trim(string(fileData), "\r\n"))
					req.Header.Set("Connection", "close")

					if req != nil {
						defer req.Body.Close()
						body, _ = ioutil.ReadAll(req.Body)
					}

					if err != nil {
						utils.PrintError("Redis import", "Can't create request to API to recieve flows", tdsModuleName)
					}

					if req.Status == "200 OK" {
						if count, err := ImportFlowsToRedis(body); err != false {
							config.TDSStatistic.AppendedFlows += count
							config.TDSStatistic.UpdatedFlows++

							config.Telegram.SendMessage("\n" + timestampPrintable + "\n" +
								config.Cfg.General.Name + "\nRequested flows from API\n" +
								"\nUpdated flows: " + strconv.Itoa(count) +
								"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))

							utils.PrintInfo("Redis import", "updated flows successful", tdsModuleName)

							// saving current timestamp to file
							ioutil.WriteFile(timestampFile, []byte(timestampWriteable), 0644)

						} else {
							// config.Telegram.SendMessage("```\n" + timestampPrintable + "\n" +
							// 	config.Cfg.General.Name + "\nRequested flows from API\n" +
							// 	"\nUpdated flows: nothing to update" +
							// 	"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG) +
							// 	"```")
						}

					} else {
						utils.PrintDebug("Error", "Recieving new flows failed", tdsModuleName)
						//
						// config.Telegram.SendMessage("```\n" + timestampPrintable + "\n" +
						// 	config.Cfg.General.Name + "\nRequested flows from API\n" +
						// 	"\nError: can't connect to API service" +
						// 	"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG) +
						// 	"```")
					}
				} else {
					// ---------------------------------------------------------------------------------------------------------------------
					// DEFAULT LOADING
					// When TDS starting up first time we need to load all flows in it
					// ---------------------------------------------------------------------------------------------------------------------
					fileData = []byte(defaultStartOfEpoch) // 2000-01-01

					// performing request to our API
					url := config.Cfg.Redis.ApiFlowsURL
					req, err := http.Get(url + strings.Trim(string(fileData), "\r\n"))

					if req != nil {
						defer req.Body.Close()
						body, _ = ioutil.ReadAll(req.Body)
					}

					if err != nil {
						utils.PrintError("Redis import", "Can't create request to API to recieve flows", tdsModuleName)
					}

					if req.Status == "200 OK" {
						if count, err := ImportFlowsToRedis(body); err != false {
							config.TDSStatistic.UpdatedFlows++
							// writing debug
							config.Telegram.SendMessage("\n" + timestampPrintable + "\n" +
								config.Cfg.General.Name + "\nRequested flows from API\n" +
								"\nUpdated flows: " + strconv.Itoa(count) +
								"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))

							utils.PrintInfo("Redis import", "All flows loaded successful", tdsModuleName)

							// saving current timestamp to file
							ioutil.WriteFile(timestampFile, []byte(timestampWriteable), 0644)

						} else {
							// config.Telegram.SendMessage("```\n" + timestampPrintable + "\n" +
							// 	config.Cfg.General.Name + "\nRequested flows from API\n" +
							// 	"\nUpdated flows: 0 nothing to update" +
							// 	"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG) +
							// 	"```")
						}

					} else {
						utils.PrintDebug("Error", "Recieving new flows failed", tdsModuleName)
						// config.Telegram.SendMessage("```\n" + timestampPrintable + "\n" +
						// 	config.Cfg.General.Name + "\nRequested flows from API\n" +
						// 	"\nUpdated flows: 0" +
						// 	"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG) +
						// 	"```")
					}
				}

				//defer runtime.GC() // startup garbage collector
				time.Sleep(time.Duration(1+config.Cfg.Redis.UpdateFlows) * time.Minute)
			}
		}
	}()
	return c
}
