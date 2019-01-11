/****************************************************************************************************
*
* Sending to API channel or saving to file
* special for Meta CPA, Ltd.
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
	"github.com/predatorpc/durafmt"
	"io"
	"metatds/utils"
	"net/http"
	"os"
	"strconv"
	"time"
)

const sendModuelName = "send.go"

func RedisSendOrSaveClicks() <-chan string {
	c := make(chan string)

	go func() {
		for {
			var clicks []map[string]string
			var KeysToDelete []string

			if IsRedisAlive {
				t := time.Now()

				var value int64
				if Cfg.Click.MaxDropItems > 0 {
					value = int64(Cfg.Click.MaxDropItems)
				} else {
					value = 10000
				}

				keys, _, _ := Redisdb.Scan(0, "*:click:*", value).Result()

				for _, item := range keys {
					KeysToDelete = append(KeysToDelete, item)
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

					if err != nil {
						recover()
					} else {
						// we shouldn't write to request if doesn't created
						req.Header.Set("X-Token", token)
						req.Header.Set("Content-Type", "application/json")
						req.Header.Set("Connection", "close")
					}

					client := &http.Client{}
					resp, err := client.Do(req)

					if err != nil {
						recover()
						if Cfg.Debug.Level > 0 {
							utils.PrintError("API Error", "Can`t send clicks to\n URL = "+url, sendModuelName)
						}
					} else {

						if resp != nil {

							utils.PrintInfo("Response status", resp.Status, sendModuelName)

							if resp.Status == "200 OK" {
								body := bytes.NewBuffer(nil)
								_, _ = io.Copy(body, resp.Body)
								_ = resp.Body.Close()

								utils.PrintError("Response", body.String(), sendModuelName)

								if KeysToDelete != nil {
									_ = Redisdb.MDel(KeysToDelete).Err()
								}

								if Cfg.Debug.Level > 1 {
									Telegram.SendMessage("\n" + utils.CURRENT_TIMESTAMP + "\n" +
										Cfg.General.Name + "\nClicks sent to API: " + strconv.Itoa(TDSStatistic.ClicksSentToRedis) +
										"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))
								}
							} else {
								body := bytes.NewBuffer(nil)
								_, _ = io.Copy(body, resp.Body)
								_ = resp.Body.Close()

								utils.PrintError("Error response", body.String(), sendModuelName)

								utils.CreateDirIfNotExist("clicks")
								f, _ := os.OpenFile("clicks/"+utils.CURRENT_TIMESTAMP_FS+".json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
								_, _ = f.WriteString(string(jsonData))
								_ = f.Close()

								if KeysToDelete != nil {
									_ = Redisdb.MDel(KeysToDelete).Err()
								}

								Telegram.SendMessage("\n" + utils.CURRENT_TIMESTAMP + "\n" +
									Cfg.General.Name + "\nClicks saved to file" +
									"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))
							}
						} else {
							utils.PrintError("Error response", "Can't read response...", sendModuelName)
						}
					}
				}

				if Cfg.Debug.Level > 1 {
					fmt.Println("Time elapsed total: ", time.Since(t))
				}
				jsonData = nil
			}

			KeysToDelete = nil
			clicks = nil
			time.Sleep(time.Duration(1+Cfg.Click.DropToRedis) * time.Second)
		}
	}()
	return c
}
