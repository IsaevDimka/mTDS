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

			// Если Redis жив, то мы работаем, если нет то следующая итерация
			if IsRedisAlive {
				t := time.Now() // засекаем время

				// если у нас нет текущей настройки в конфиге берем 10000 по умолчанию
				var value int64
				if Cfg.Click.MaxDropItems > 0 {
					value = int64(Cfg.Click.MaxDropItems)
				} else {
					value = 10000
				}

				// Получаем сканом т.к. быстрее чем кейс
				keys, _, _ := Redisdb.Scan(0, "*:click:*", value).Result()

				// Получаем наши данные в массив структур кликов
				// Сразу решаем какие будем удалять
				for _, item := range keys {
					KeysToDelete = append(KeysToDelete, item)
					d, _ := Redisdb.HGetAll(item).Result()
					clicks = append(clicks, d)
				}

				// Собираем гигантский JSON
				jsonData, _ := json.Marshal(clicks)

				if Cfg.Debug.Level > 1 {
					fmt.Println("Time elapsed export: ", time.Since(t))
				}

				// Если данные состряпались, тогда мы их отправим
				if len(jsonData) > 0 && len(clicks) > 0 {

					url := Cfg.Click.ApiUrl
					token := Cfg.Click.ApiToken
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
							// пишем лог
							utils.PrintInfo("Response status", resp.Status, sendModuelName)

							if resp.Status == "200 OK" {
								// отмечаем, что мы их отправили
								TDSStatistic.ClicksSentToRedis += len(clicks)

								// READING RESPONSE (NECCESSARY) ------------------------------------------------------------------------------
								// Если не будем читать ответы память будем кушать over dahua
								body := bytes.NewBuffer(nil)
								// Функция copy в отличие от функции read читает буфером по 32 Кб т.е использование памяти никакое
								// Важно!!!
								_, _ = io.Copy(body, resp.Body)
								_ = resp.Body.Close()
								// ------------------------------------------------------------------------------------------------------------

								utils.PrintError("Response", body.String(), sendModuelName)

								// Если мы хотим подебужить что мы там отправляем нужна настройка в сеттингс
								if Cfg.Click.BackupClicks {
									utils.CreateDirIfNotExist("backup")
									f, _ := os.OpenFile("backup/"+utils.CURRENT_TIMESTAMP_FS+".json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
									_, _ = f.WriteString(string(jsonData))
									_ = f.Close()
								}

								// Если у нас есть, чего удалить, то стираем все это из редиса
								// Потом незабываем обнулить все это дело
								if KeysToDelete != nil {
									_ = Redisdb.MDel(KeysToDelete).Err()
								}

								// Шлем себе дебуг мессейдж, на тему как хорошо все прошло
								if Cfg.Debug.Level > 1 {
									Telegram.SendMessage("\n" + utils.CURRENT_TIMESTAMP + "\n" +
										Cfg.General.Name + "\nClicks sent to API: " + strconv.Itoa(TDSStatistic.ClicksSentToRedis) +
										"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))
								}
							} else {
								// Если АПИ недоступно, или что-то пошло не так, то нам надо забэкапить все это дело в файло
								// сначала прочитаем также как и в коде выше
								body := bytes.NewBuffer(nil)
								_, _ = io.Copy(body, resp.Body)
								_ = resp.Body.Close()

								utils.PrintError("Error response", body.String(), sendModuelName)

								// Создадим если нет
								utils.CreateDirIfNotExist("clicks")
								// тут тоже обязательно делать через системные вызовы иначе память ку-ку
								f, _ := os.OpenFile("clicks/"+utils.CURRENT_TIMESTAMP_FS+".json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
								_, _ = f.WriteString(string(jsonData))
								_ = f.Close()

								// А Редис нам освобождать по любому
								if KeysToDelete != nil {
									_ = Redisdb.MDel(KeysToDelete).Err()
								}

								// Шлем дебаг в независимости от уровня
								// т.к. надо срочно обратить внимание на все это дело
								//
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
				// Обнуляем данные, чтобы ничего у нас не осталось вдруг, от предыдушего раза
				jsonData = nil
			}

			// Обнуляем кандидатов на удаление т.к. мы их уже удалили
			KeysToDelete = nil
			// Обнуляем клики они нам тоже уже не нужны (экномим память)
			clicks = nil
			// Поспим чуть-чуть как настроим
			time.Sleep(time.Duration(1+Cfg.Click.DropToRedis) * time.Second)
		}
	}()
	return c
}
