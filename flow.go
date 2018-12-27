package main

import (
	"fmt"
	"github.com/labstack/echo"
	"math/rand"
	"metatds/config"
	"metatds/utils"
	"net/http"
	"runtime"
	"strings"
	"time"
)

/*
*
* Get Flow information in JSON ir Redirect to some selected location
* and Redirect handler all in one
*
*
 */

func flowHandler(c echo.Context) error {

	// for Debug don't delete for a few time

	// var memory runtime.MemStats
	// runtime.ReadMemStats(&memory)
	// fmt.Print("[MEMORY USAGE]", " Alloc: ",utils.BToMb(memory.Alloc)," Mb, Total: ",utils.BToMb(memory.StackSys)," Mb, Sys: ",utils.BToMb(memory.Sys)," Mb")
	// s:=strconv.FormatUint(utils.BToMb(memory.Sys+memory.Alloc),10)
	// fmt.Println("S = ",s)

	CID, cookieError := c.Cookie("CID")

	if cookieError != nil {
		if config.Cfg.Debug.Level > 1 {
			utils.PrintDebug("Cookie", "Error reading cookie", tdsModuleName)
		}
	} else {
		if config.Cfg.Debug.Level > 1 {
			utils.PrintInfo("Cookie", "CID = "+CID.Value, tdsModuleName)
		}
	}

	if config.IsRedisAlive {

		start := time.Now()

		var Info InfoData
		var LandingTemplate, PrelandingTemplate string
		var LandingTemplateID, PrelandingTemplateID int

		resultMap := utils.URIByMap(c, keyMap) // вот в этот массив

		if cookieError == nil {
			Info.Click.Hash = CID.Value
			config.TDSStatistic.CookieRequest++
		} else {
			Info.Click.Hash = Info.Click.GenerateCID()
		}

		// генерим СИД
		resultMap["click_hash"] = append(resultMap["click_hash"], Info.Click.Hash) // запишем сразу в наш массив
		resultMap["click_id"] = append(resultMap["click_id"], Info.Click.Hash)     // support for old version TDS
		Info.Flow = Info.Flow.GetInfo(strings.Join(resultMap["flow_hash"], ""))    // получить всю инфу о потоке

		t := time.Now() // собираем данные для сейва в базу

		Info.Click.UserAgent = c.Request().UserAgent()
		Info.Click.Time = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
		Info.Click.URL = "http://" + config.Cfg.General.Host + c.Request().RequestURI

		Info.Click.IP = c.Request().RemoteAddr
		Info.Click.Referer = c.Request().Referer()

		Info.Click.FlowHash = Info.Flow.Hash
		Info.Click.FlowID = Info.Flow.ID
		Info.Click.WebMasterID = Info.Flow.WebMasterID
		Info.Click.WebMasterCurrencyID = Info.Flow.WebMasterCurrencyID
		Info.Click.OfferID = Info.Flow.OfferID

		Info.Click.Sub1 = strings.Join(resultMap["sub1"], "")
		Info.Click.Sub2 = strings.Join(resultMap["sub2"], "")
		Info.Click.Sub3 = strings.Join(resultMap["sub3"], "")
		Info.Click.Sub4 = strings.Join(resultMap["sub4"], "")
		Info.Click.Sub5 = strings.Join(resultMap["sub5"], "")

		// если есть поток и есть клик, значит можно и дальше идти
		if strings.Join(resultMap["flow_hash"], "") != "" && strings.Join(resultMap["click_hash"], "") != "" &&
			Info.Flow.ID > 0 {

			if len(Info.Flow.Lands) > 0 {
				Random := rand.Intn(len(Info.Flow.Lands))
				LandingTemplate = Info.Flow.Lands[Random].URL  // получаем рандомный урл ленда
				LandingTemplateID = Info.Flow.Lands[Random].ID // strconv.Itoa(Info.Flow.Lands[Random].ID)
			} else {
				config.TDSStatistic.IncorrectRequest++ // add counter tick
				msg := []byte(`{"code":400, "message":"No landing templates found"}`)
				return c.JSONBlob(400, msg)
			}

			if len(Info.Flow.Prelands) > 0 {
				Random := rand.Intn(len(Info.Flow.Prelands))
				PrelandingTemplate = Info.Flow.Prelands[Random].URL  // получаем рандомный урл преленда
				PrelandingTemplateID = Info.Flow.Prelands[Random].ID // strconv.Itoa(Info.Flow.Prelands[Random].ID)
			}

			// Если дебаг то печатаем все это добро
			if config.Cfg.Debug.Level > 1 && len(Info.Flow.Prelands) > 0 && len(Info.Flow.Lands) > 0 {
				utils.PrintDebug("Land", resultMap, tdsModuleName)
				utils.PrintDebug("Land", Info.Flow.Lands[rand.Intn(len(Info.Flow.Lands))].URL, tdsModuleName)
				utils.PrintDebug("Preland", Info.Flow.Prelands[rand.Intn(len(Info.Flow.Prelands))].URL, tdsModuleName)
			}

			// ставим куку на этот урл
			cookie := new(http.Cookie)
			cookie.Name = "CID"
			cookie.Value = Info.Click.Hash
			cookie.Expires = time.Now().Add(365 * 24 * time.Hour) // for an year

			// ----------------------------------------------------------------------------------------------------
			// выбираем по формату, что будем отдавать
			// редирект на лендинг
			// LAND
			// ----------------------------------------------------------------------------------------------------
			if strings.Join(resultMap["format"], "") == "land" {
				for _, item := range keyMap {
					LandingTemplate = strings.Replace(LandingTemplate, fmt.Sprintf("{%s}", item),
						strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
				}

				if config.Cfg.Debug.Level > 1 {
					utils.PrintInfo("Result", LandingTemplate, tdsModuleName)
				}

				Info.Click.LandingID = LandingTemplateID
				Info.Click.PrelandingID = 0
				Info.Click.IsVisitedLP = 1
				Info.Click.LocationLP = LandingTemplate

				cookie.Path = Info.Click.LocationLP
				c.SetCookie(cookie)

				defer Info.Click.Save()

				config.TDSStatistic.RedirectRequest++ // add counter tick
				config.TDSStatistic.ProcessingTime += time.Since(start)

				if config.Cfg.Debug.Level > 1 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				if len(config.ResponseAverage) <= 100 {
					utils.ResponseAverage = append(utils.ResponseAverage, time.Since(start))
				} else {
					utils.ResponseAverage = nil
				}

				if !config.Cfg.Debug.Test {
					runtime.GC()
					return c.Redirect(302, LandingTemplate)
				} else {
					runtime.GC()
					return c.Blob(200, "image/png", pixel)
				}
			}
			// ----------------------------------------------------------------------------------------------------
			// PRELAND
			// редирект на пре-лендинг
			// ----------------------------------------------------------------------------------------------------
			if strings.Join(resultMap["format"], "") == "preland" {
				if len(Info.Flow.Prelands) <= 0 {
					config.TDSStatistic.IncorrectRequest++ // add counter tick

					msg := []byte(`{"code":400, "message":"No pre-landing templates found"}`)
					return c.JSONBlob(400, msg)
				}

				for _, item := range keyMap {
					PrelandingTemplate = strings.Replace(PrelandingTemplate, fmt.Sprintf("{%s}", item),
						strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
				}

				if config.Cfg.Debug.Level > 1 {
					utils.PrintInfo("Result", PrelandingTemplate, tdsModuleName)
				}

				Info.Click.LandingID = 0
				Info.Click.IsVisitedPL = 1
				Info.Click.PrelandingID = PrelandingTemplateID
				Info.Click.LocationPL = PrelandingTemplate

				cookie.Path = Info.Click.LocationPL
				c.SetCookie(cookie)

				defer Info.Click.Save()

				config.TDSStatistic.RedirectRequest++ // add counter tick
				config.TDSStatistic.ProcessingTime += time.Since(start)

				if config.Cfg.Debug.Level > 1 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				if len(config.ResponseAverage) <= 100 {
					utils.ResponseAverage = append(utils.ResponseAverage, time.Since(start))
				} else {
					utils.ResponseAverage = nil
				}

				if !config.Cfg.Debug.Test {

					runtime.GC()
					return c.Redirect(302, PrelandingTemplate)
				} else {

					runtime.GC()
					return c.Blob(200, "image/png", pixel)
				}
			}
			// ----------------------------------------------------------------------------------------------------
			// JSON FORMAT
			// отдать данные потока в джейсоне красиво
			// ----------------------------------------------------------------------------------------------------
			if strings.Join(resultMap["format"], "") == "json" {
				if len(Info.Flow.Prelands) > 0 {
					for _, item := range keyMap {
						PrelandingTemplate = strings.Replace(PrelandingTemplate, fmt.Sprintf("{%s}", item),
							strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
					}
					Info.Flow.RandomPreland = PrelandingTemplate
				}

				if len(Info.Flow.Lands) > 0 {
					for _, item := range keyMap {
						LandingTemplate = strings.Replace(LandingTemplate, fmt.Sprintf("{%s}", item),
							strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
					}
					Info.Flow.RandomLand = LandingTemplate
				}

				Info.Click.LandingID = LandingTemplateID
				Info.Click.LocationLP = LandingTemplate
				Info.Click.IsVisitedLP = 0

				Info.Click.PrelandingID = PrelandingTemplateID
				Info.Click.LocationPL = PrelandingTemplate
				Info.Click.IsVisitedPL = 1

				cookie.Path = Info.Click.LocationLP
				c.SetCookie(cookie)

				s := utils.JSONPretty(Info)

				defer Info.Click.Save()

				if config.Cfg.Debug.Level > 1 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				if len(config.ResponseAverage) <= 100 {
					utils.ResponseAverage = append(utils.ResponseAverage, time.Since(start))
				} else {
					utils.ResponseAverage = nil
				}

				config.TDSStatistic.FlowInfoRequest++ // add counter tick
				config.TDSStatistic.ProcessingTime += time.Since(start)

				runtime.GC()
				return c.String(200, s)

			} else {
				// ----------------------------------------------------------------------------------------------------
				// OLD School
				// ----------------------------------------------------------------------------------------------------
				// если никаких ключей нет, то пробрасываем дальше (по старой схеме)
				// на первый выбраный домен из списка если несколько или на первый
				// ----------------------------------------------------------------------------------------------------
				for _, item := range keyMap {
					LandingTemplate = strings.Replace(LandingTemplate, fmt.Sprintf("{%s}", item),
						strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
				}

				Info.Click.LandingID = LandingTemplateID
				Info.Click.PrelandingID = 0
				Info.Click.LocationLP = LandingTemplate

				cookie.Path = Info.Click.LocationLP
				c.SetCookie(cookie)

				defer Info.Click.Save()
				config.TDSStatistic.RedirectRequest++ // add counter tick
				config.TDSStatistic.ProcessingTime += time.Since(start)

				if config.Cfg.Debug.Level > 1 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				if len(config.ResponseAverage) <= 100 {
					utils.ResponseAverage = append(utils.ResponseAverage, time.Since(start))
				} else {
					utils.ResponseAverage = nil
				}

				if !config.Cfg.Debug.Test {

					runtime.GC()

					return c.Redirect(302, LandingTemplate)
				} else {

					runtime.GC()
					return c.Blob(200, "image/png", pixel)
				}
			}
		} else {
			config.TDSStatistic.IncorrectRequest++ // add counter tick
			// если нет клика или потока, то все привет
			msg := []byte(`{"code":400, "message":"Insuficient parameters supplied"}`)

			runtime.GC()
			return c.JSONBlob(400, msg)
		}
	} else {

		// если нет редиски, то все привет
		msg := []byte(`{"code":500, "message":"No connection to RedisDB"}`)

		runtime.GC()
		return c.JSONBlob(400, msg)
	}
}
