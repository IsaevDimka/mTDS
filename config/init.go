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
	"net/http"
	"time"
)

const initModuleName = "init.go"

var UpTime time.Time

//
// Main initialization handle
//
func init() {
	// get current timestamp
	UpTime = time.Now()
	utils.CURRENT_TIMESTAMP = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		UpTime.Year(), UpTime.Month(), UpTime.Day(), UpTime.Hour(), UpTime.Minute(), UpTime.Second())

	// поставищик текущей метки времени и ее строкового представления
	CurrentTimeStampTicker()

	// сначала загружаем настройки потом, цепляем все остальное
	InitConfig()

	if Cfg.Debug.Level > 0 {
		utils.PrintDebug("Initialization", "", initModuleName)
	}

	// issue with too many open files
	http.DefaultClient.Timeout = time.Second * time.Duration(1+Cfg.General.HTTPTimeout)

	// цепляем редис и потом, проверяем постоянно, как у него дела
	RedisDBChan()

	// Напишем всем, что мы стартанули
	tlgrmRecipients := utils.Explode(Cfg.Telegram.Recipients, "; ")
	tlgrm := Telegram.Init(tlgrmRecipients, Cfg.Telegram.Socks5User, Cfg.Telegram.Socks5Password,
		Cfg.Telegram.Socks5Proxy, Cfg.Telegram.ApiURL, Cfg.Telegram.Token, Cfg.Telegram.UseProxy)

	// timeStamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
	// 	UpTime.Year(), UpTime.Month(), UpTime.Day(), UpTime.Hour(), UpTime.Minute(), UpTime.Second())

	Telegram.SendMessage("\n" + utils.CURRENT_TIMESTAMP + "\n" + Cfg.General.Name + "\nTDS Service started\n")

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

//
// Telegram send statistic channel
//
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
			time.Sleep(time.Duration(1+Cfg.Telegram.MsgInterval) * time.Second) // поспим чуть чуть
		}
	}()

	return c
}

//
// Reload config channel
//
func ReloadConfigChan() <-chan string {
	c := make(chan string)
	go func() {
		for {
			// перезагружаем конфиг и идем спать
			ReloadConfig()
			// поспим чуть чуть +1 its to avoid dumbs with zero multiplication
			time.Sleep(time.Duration(1+Cfg.General.ConfReload) * time.Second)
		}
	}()
	return c
}

//
// TimeStamp ticker
//
func CurrentTimeStampTicker() <-chan string {
	c := make(chan string)
	go func() {
		for {
			t := time.Now()
			utils.CURRENT_UNIXTIME = t
			utils.CURRENT_TIMESTAMP = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
				t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
			utils.CURRENT_TIMESTAMP_FS = fmt.Sprintf("%d%02d%02d%02d%02d%02d",
				t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

			time.Sleep(time.Second * 1)
		}
	}()
	return c
}
