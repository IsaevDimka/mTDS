/****************************************************************************************************
*
* Flow import module, special for Meta CPA, Ltd.
* by Michael S. Merzlyakov AFKA predator_pc@09012019
* version v2.0.5
*
* created at 04122018
* last edit: 09012019
*
*****************************************************************************************************/

package main

import (
	"bytes"
	"encoding/json"
	"github.com/predatorpc/durafmt"
	"io"
	"io/ioutil"
	"metatds/config"
	"metatds/models"
	"metatds/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const importModuleName = "import.go"

func ImportFlowsToRedis(jsonData []byte) (int, bool) {
	Flows := make(map[string]models.FlowImportData)
	if err := json.Unmarshal(jsonData, &Flows); err != nil {
		if config.Cfg.Debug.Level > 1 {
			utils.PrintDebug("Error", "Can`t decode JSON given", importModuleName)
		}
	} else {
		for _, item := range Flows {
			params := make(map[string]interface{})
			params["ID"] = item.ID
			params["Hash"] = item.Hash
			params["OfferID"] = item.OfferID
			params["WebMasterID"] = item.WebMasterID
			params["WebMasterCurrencyID"] = item.WebMasterCurrencyID

			_ = config.Redisdb.HMSet(item.Hash, params).Err()

			if len(item.Lands) > 0 {
				for _, lands := range item.Lands {
					_ = config.Redisdb.HSet(item.Hash+":lands", strconv.Itoa(lands.ID), lands.URL).Err()
				}
			}

			if len(item.Prelands) > 0 {
				for _, prelands := range item.Prelands {
					_ = config.Redisdb.HSet(item.Hash+":prelands", strconv.Itoa(prelands.ID), prelands.URL).Err()
				}
			}

			if len(item.Counters) > 0 {
				for _, counters := range item.Counters {
					_ = config.Redisdb.HSet(item.Hash+":counters", counters.Name, strconv.Itoa(counters.ID)).Err()
				}
			}
		}
		return len(Flows), true
	}
	return len(Flows), false
}

func UpdateFlowsListChan() <-chan string {
	c := make(chan string)
	go func() {
		for {
			if config.IsRedisAlive {
				t := time.Now() // start counting elapsed time
				body := bytes.NewBuffer(nil)
				// writing it to file
				timestampWriteable := strconv.FormatInt(time.Now().Unix(), 10)
				fileData, err := ioutil.ReadFile(timestampFile)

				if err == nil {
					// ---------------------------------------------------------------------------------------------------------------------
					// LOADING WITH PARAMS OF LAST UPDATE
					// When TDS starting up first time we need to load all flows in it
					// ---------------------------------------------------------------------------------------------------------------------
					// if we can't parse this we should get all anyway
					_, err := strconv.ParseInt(string(fileData), 10, 64)

					// getting current count of flows and if it isn't null then proceeed
					// checking this for error has no effect if Redis is Alive
					currentCount, _, _ := config.Redisdb.Scan(0, "", 10000).Result()

					// if cannot decode label setting default one
					if err != nil || len(currentCount) == 0 {
						if config.Cfg.Debug.Level > 0 {
							utils.PrintDebug("Error", "parsing timestamp from `"+timestampFile+"` failure", importModuleName)
						}
						fileData = []byte(defaultStartOfEpoch) // 2000-01-01
					}

					// performing request to our API
					url := config.Cfg.Redis.ApiFlowsURL
					req, err := http.Get(url + strings.Trim(string(fileData), "\r\n"))

					if req != nil {
						// setting header in case of API request is not NIL
						req.Header.Set("Connection", "close")
						// reading the body
						_, _ = io.Copy(body, req.Body)
						// closing anyway now
						// defer is not needed cause we get an exception before
						_ = req.Body.Close()
					}

					if err != nil && config.Cfg.Debug.Level > 0 {
						utils.PrintError("Redis import", "Can't create request to API to recieve flows: \n URL = "+
							url+strings.Trim(string(fileData), "\r\n"), importModuleName)
					} else {
						if req != nil {
							if req.Status == "200 OK" {
								count, err := ImportFlowsToRedis(body.Bytes())
								if err != false {
									config.TDSStatistic.AppendedFlows += count
									config.TDSStatistic.UpdatedFlows++

									config.Telegram.SendMessage("\n" + utils.CURRENT_TIMESTAMP + "\n" +
										config.Cfg.General.Name + "\nRequested flows from API\n" +
										"\nUpdated flows: " + strconv.Itoa(count) +
										"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))

									utils.PrintInfo("Redis import", "updated flows successful", importModuleName)
									// saving current timestamp to file
									_ = ioutil.WriteFile(timestampFile, []byte(timestampWriteable), 0644)
								} else {
									utils.PrintDebug("Error", "Normal bootstrap: Writing to Redis failed or empty response", importModuleName)
								}

							} else {
								utils.PrintDebug("Error", "Receiving new flows failed", importModuleName)
							}
						} else {
							utils.PrintDebug("Error", "1 Can't read response: Receiving new flows failed", importModuleName)
						}
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
						// setting header in case of API request is not NIL
						req.Header.Set("Connection", "close")
						// reading the body
						_, _ = io.Copy(body, req.Body)
						// closing anyway now
						// defer is not needed cause we get an exception before
						_ = req.Body.Close()
					}

					if err != nil {
						utils.PrintError("Redis import", "Can't create request to API to recieve flows: \n URL = "+
							url+strings.Trim(string(fileData), "\r\n"), importModuleName)
					} else {
						if req != nil {
							if req.Status == "200 OK" {
								count, err := ImportFlowsToRedis(body.Bytes())
								if err != false {
									config.TDSStatistic.UpdatedFlows++
									// writing debug
									config.Telegram.SendMessage("\n" + utils.CURRENT_TIMESTAMP + "\n" +
										config.Cfg.General.Name + "\nRequested flows from API\n" +
										"\nUpdated flows: " + strconv.Itoa(count) +
										"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))

									utils.PrintInfo("Redis import", "All flows loaded successful", importModuleName)
									// saving current timestamp to file
									_ = ioutil.WriteFile(timestampFile, []byte(timestampWriteable), 0644)
								} else {
									utils.PrintDebug("Error", "Writing to Redis failed", importModuleName)
								}

							} else {
								utils.PrintDebug("Error", "Emergency bootstrap: Writing to Redis failed or empty response", importModuleName)
							}
						} else {
							utils.PrintDebug("Error", "0 Can't read response: Receiving new flows failed", importModuleName)
						}
					}
				}
				body = nil
				time.Sleep(time.Duration(1+config.Cfg.Redis.UpdateFlows) * time.Second)
			}
		}
	}()
	return c
}
