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
					io.Copy(w, file)
					file.Close()

					url := Cfg.Click.ApiUrl                     // "http://116.202.27.130/set/hits"
					token := Cfg.Click.ApiToken                 // "PaILgFTQQCvX9tzS"
					req, err := http.NewRequest("POST", url, w) //bytes.NewBuffer(fileData))

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

						utils.PrintDebug("Response status", resp.Status, resendModuelName)

						if resp.Status == "200 OK" {
							body := bytes.NewBuffer(nil)
							io.Copy(body, resp.Body)
							resp.Body.Close()

							utils.PrintError("Response", body.String(), resendModuelName)
							// удаляем файл, мы его успешно обработали
							os.Remove(item)

							Telegram.SendMessage("\n" + utils.CURRENT_TIMESTAMP + "\n" +
								Cfg.General.Name + "\nResending file succedeed " + fdsReplace + " to API" +
								"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))
						} else {
							body := bytes.NewBuffer(nil)
							io.Copy(body, resp.Body)
							resp.Body.Close()

							utils.PrintInfo("Error response", body.String(), resendModuelName)
							utils.PrintDebug("Error", "Sending file to click API failed", resendModuelName)

							Telegram.SendMessage("\n" + utils.CURRENT_TIMESTAMP + "\n" +
								Cfg.General.Name + "\nResending file failed " + fdsReplace + " to API" +
								"\nTime elsapsed for operation: " + durafmt.Parse(time.Since(t)).String(durafmt.DF_LONG))
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
