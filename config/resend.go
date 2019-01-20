/****************************************************************************************************
*
* Sending to saved files .json to API
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
	"github.com/predatorpc/durafmt"
	"io"
	"metatds/utils"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const resendModuelName = "resend.go"

//
// Send -*.json stored in files to reciever API
//
func SendFileToRecieveApi() <-chan string {
	c := make(chan string)
	go func() {
		for {
			var fdsReplace string
			t := time.Now()
			fds, _ := filepath.Glob("clicks/*.json")
			if len(fds) > 0 {
				for _, item := range fds {
					fdsReplace = filepath.Base(item)

					// прочитываем весь файл в буфер по 32кб
					file, _ := os.Open(item)
					w := bytes.NewBuffer(nil)
					_, _ = io.Copy(w, file)
					_ = file.Close()

					url := Cfg.Click.ApiUrl
					token := Cfg.Click.ApiToken
					req, err := http.NewRequest("POST", url, w)

					if err != nil {
						recover()
					} else {
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
							utils.PrintDebug("Response status", resp.Status, resendModuelName)

							if resp.Status == "200 OK" {
								body := bytes.NewBuffer(nil)
								_, _ = io.Copy(body, resp.Body)
								_ = resp.Body.Close()

								utils.PrintError("Response", body.String(), resendModuelName)
								// удаляем файл, мы его успешно обработали
								_ = os.Remove(item)

								Telegram.SendMessage("\n" + utils.CURRENT_TIMESTAMP + "\n" +
									Cfg.General.Name + "\nResending file succedeed " + fdsReplace + " to API" +
									"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))
							} else {
								body := bytes.NewBuffer(nil)
								_, _ = io.Copy(body, resp.Body)
								_ = resp.Body.Close()

								utils.PrintInfo("Error response", body.String(), resendModuelName)
								utils.PrintDebug("Error", "Sending file to click API failed", resendModuelName)

								Telegram.SendMessage("\n" + utils.CURRENT_TIMESTAMP + "\n" +
									Cfg.General.Name + "\nResending file failed " + fdsReplace + " to API" +
									"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))
							}
						} else {
							utils.PrintDebug("Error", "0 Can't read response: Sending clicks failed", resendModuelName)
						}
					}

					w = nil
					// поспим между файлами
					time.Sleep(time.Second * 1)
				}
			}

			fds = nil
			time.Sleep(time.Duration(1+Cfg.Click.DropFilesToAPI) * time.Second)
		}
	}()
	return c
}
