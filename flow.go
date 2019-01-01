package main

import (
	"fmt"
	"github.com/labstack/echo"
	"math/rand"
	"metatds/config"
	"metatds/utils"
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
 const minimumStatCount = 1000

func flowHandler(c echo.Context) error {
	if config.IsRedisAlive {
		// начинаем замер производительности
		start := time.Now()
		var Info InfoData
		var LandingTemplate, PrelandingTemplate string
		var LandingTemplateID, PrelandingTemplateID int
		//------------------------------------------------------------------------------------------------------
		// Читаем параметры с которыми мы хотим редирект
		//------------------------------------------------------------------------------------------------------
		resultMap := utils.URIByMap(c, keyMap) // вот в этот массив
		//------------------------------------------------------------------------------------------------------
		// Читаем куку
		//------------------------------------------------------------------------------------------------------
		CID, cookieError := c.Cookie("CID")
		if cookieError == nil {
			if config.Cfg.Debug.Level > 1 {
				utils.PrintInfo("Cookie", "CID = "+CID.Value, tdsModuleName)
			}
			Info.Click.Hash = CID.Value
			config.TDSStatistic.CookieRequest++
		} else {
			if config.Cfg.Debug.Level > 1 {
				utils.PrintDebug("Cookie", "Error reading cookie", tdsModuleName)
			}
			// генерим СИД
			Info.Click.Hash = Info.Click.GenerateCID()
		}
		//------------------------------------------------------------------------------------------------------

		resultMap["click_hash"] = append(resultMap["click_hash"], Info.Click.Hash) // запишем сразу в наш массив
		resultMap["click_id"] = append(resultMap["click_id"], Info.Click.Hash)     // support for old version TDS

		//------------------------------------------------------------------------------------------------------
		// Тут вот может быть можно ускорить
		//------------------------------------------------------------------------------------------------------
		Info.Flow = Info.Flow.GetInfo(strings.Join(resultMap["flow_hash"], "")) // получить всю инфу о потоке
		//------------------------------------------------------------------------------------------------------
		Info.Click.Time = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute(), start.Second())

		// _SERVER information
		Info.Click.UserAgent = c.Request().UserAgent()
		Info.Click.IP = c.Request().RemoteAddr
		Info.Click.Referer = c.Request().Referer()

		// грузим в клик все из потока
		Info.Click.FlowHash = Info.Flow.Hash
		Info.Click.FlowID = Info.Flow.ID
		Info.Click.WebMasterID = Info.Flow.WebMasterID
		Info.Click.WebMasterCurrencyID = Info.Flow.WebMasterCurrencyID
		Info.Click.OfferID = Info.Flow.OfferID

		// грузим субаки
		Info.Click.URL = "http://" + config.Cfg.General.Host + c.Request().RequestURI
		Info.Click.Sub1 = strings.Join(resultMap["sub1"], "")
		Info.Click.Sub2 = strings.Join(resultMap["sub2"], "")
		Info.Click.Sub3 = strings.Join(resultMap["sub3"], "")
		Info.Click.Sub4 = strings.Join(resultMap["sub4"], "")
		Info.Click.Sub5 = strings.Join(resultMap["sub5"], "")

		// если есть поток и есть клик, значит можно и дальше идти
		//------------------------------------------------------------------------------------------------------
		// Подготовка всей инфы
		//------------------------------------------------------------------------------------------------------
		if strings.Join(resultMap["flow_hash"], "") != "" && strings.Join(resultMap["click_hash"], "") != "" && Info.Flow.ID > 0 {

			//------------------------------------------------------------------------------------------------------
			// Выбор лендингов и прелендингов куда буем редиректить, если лендов нет, то вообще заканчиваем цирк
			//------------------------------------------------------------------------------------------------------
			if len(Info.Flow.Lands) > 0 {
				Random := rand.Intn(len(Info.Flow.Lands))
				LandingTemplate = Info.Flow.Lands[Random].URL // получаем рандомный урл ленда
				LandingTemplateID = Info.Flow.Lands[Random].ID

				for _, item := range keyMap {
					LandingTemplate = strings.Replace(LandingTemplate, fmt.Sprintf("{%s}", item),
						strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
				}
				Info.Flow.RandomLand = LandingTemplate

			} else {
				config.TDSStatistic.IncorrectRequest++ // add counter tick
				msg := []byte(`{"code":400, "message":"No landing templates found"}`)
				return c.JSONBlob(400, msg)
			}

			if len(Info.Flow.Prelands) > 0 {
				Random := rand.Intn(len(Info.Flow.Prelands))
				PrelandingTemplate = Info.Flow.Prelands[Random].URL // получаем рандомный урл преленда
				PrelandingTemplateID = Info.Flow.Prelands[Random].ID

				for _, item := range keyMap {
					PrelandingTemplate = strings.Replace(PrelandingTemplate, fmt.Sprintf("{%s}", item),
						strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
				}
				Info.Flow.RandomPreland = PrelandingTemplate
			}
			//------------------------------------------------------------------------------------------------------

			// Если дебаг то печатаем все это добро
			if config.Cfg.Debug.Level > 1 && len(Info.Flow.Prelands) > 0 && len(Info.Flow.Lands) > 0 {
				utils.PrintDebug("Land", resultMap, tdsModuleName)
				utils.PrintDebug("Land", Info.Flow.Lands[rand.Intn(len(Info.Flow.Lands))].URL, tdsModuleName)
				utils.PrintDebug("Preland", Info.Flow.Prelands[rand.Intn(len(Info.Flow.Prelands))].URL, tdsModuleName)
			}

			// ----------------------------------------------------------------------------------------------------
			// выбираем по формату, что будем отдавать
			// редирект на лендинг
			// LAND
			// ----------------------------------------------------------------------------------------------------
			if strings.Join(resultMap["format"], "") == "land" {
				// добиваем клик нужной инфой, теперь можем его записывать
				Info.Click.LandingID = LandingTemplateID
				Info.Click.PrelandingID = 0
				Info.Click.IsVisitedLP = 1
				Info.Click.LocationLP = LandingTemplate

				// в отдельной горутине запускаем сохранение в редис
				// не останавливая основной поток!
				// и куку тудаже если надо ее поставить
				defer Info.Click.Save()
				if cookieError != nil {
					defer c.SetCookie(utils.SaveCookieToUser(Info.Click.Hash, Info.Click.LocationLP))
				}

				//------------------------------------------------------------------------------------------------------
				// добавляем статистику
				//------------------------------------------------------------------------------------------------------
				config.TDSStatistic.RedirectRequest++ // add counter tick
				config.TDSStatistic.ProcessingTime += time.Since(start)

				if config.Cfg.Debug.Level > 1 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				if len(utils.ResponseAverage) < minimumStatCount {
					utils.ResponseAverage = append(utils.ResponseAverage, time.Since(start))
				} else {
					utils.ResponseAverage = utils.ResponseAverageDefault
				}

				//------------------------------------------------------------------------------------------------------
				// Финал редиректим
				//------------------------------------------------------------------------------------------------------//------------------------------------------------------------------------------------------------------
				if !config.Cfg.Debug.Test {
				//	defer runtime.GC()
					return c.Redirect(302, LandingTemplate)
				} else {
				//	defer runtime.GC()
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

				Info.Click.LandingID = 0
				Info.Click.IsVisitedPL = 1
				Info.Click.PrelandingID = PrelandingTemplateID
				Info.Click.LocationPL = PrelandingTemplate

				defer Info.Click.Save()
				if cookieError != nil {
					defer c.SetCookie(utils.SaveCookieToUser(Info.Click.Hash, Info.Click.LocationLP))
				}

				// ----------------------------------------------------------------------------------------------------
				// STATS
				// ----------------------------------------------------------------------------------------------------
				config.TDSStatistic.RedirectRequest++ // add counter tick
				config.TDSStatistic.ProcessingTime += time.Since(start)

				if config.Cfg.Debug.Level > 1 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				if len(utils.ResponseAverage) < minimumStatCount {
					utils.ResponseAverage = append(utils.ResponseAverage, time.Since(start))
				} else {
					utils.ResponseAverage = utils.ResponseAverageDefault
				}

				// ----------------------------------------------------------------------------------------------------
				// FINAL
				// ----------------------------------------------------------------------------------------------------
				if !config.Cfg.Debug.Test {
				//	defer runtime.GC()
					return c.Redirect(302, PrelandingTemplate)
				} else {
				//	defer runtime.GC()
					return c.Blob(200, "image/png", pixel)
				}
			}
			// ----------------------------------------------------------------------------------------------------
			// JSON FORMAT
			// отдать данные потока в джейсоне красиво
			// ----------------------------------------------------------------------------------------------------
			if strings.Join(resultMap["format"], "") == "json" {
				Info.Click.LandingID = LandingTemplateID
				Info.Click.LocationLP = LandingTemplate
				Info.Click.IsVisitedLP = 0
				Info.Click.PrelandingID = PrelandingTemplateID
				Info.Click.LocationPL = PrelandingTemplate
				Info.Click.IsVisitedPL = 1

				defer Info.Click.Save()
				if cookieError != nil {
					defer c.SetCookie(utils.SaveCookieToUser(Info.Click.Hash, Info.Click.LocationLP))
				}

				// ----------------------------------------------------------------------------------------------------
				// STATS
				// ----------------------------------------------------------------------------------------------------
				config.TDSStatistic.FlowInfoRequest++ // add counter tick
				config.TDSStatistic.ProcessingTime += time.Since(start)

				if config.Cfg.Debug.Level > 1 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				if len(utils.ResponseAverage) < minimumStatCount {
					utils.ResponseAverage = append(utils.ResponseAverage, time.Since(start))
				} else {
					utils.ResponseAverage = utils.ResponseAverageDefault
				}

				// ----------------------------------------------------------------------------------------------------
				// FINAL
				// ----------------------------------------------------------------------------------------------------
				s := utils.JSONPretty(Info)
				//defer runtime.GC()
				return c.String(200, s)
			} else {
				// ----------------------------------------------------------------------------------------------------
				// OLD School
				// ----------------------------------------------------------------------------------------------------
				// если никаких ключей нет, то пробрасываем дальше (по старой схеме)
				// на первый выбраный домен из списка если несколько или на первый
				// ----------------------------------------------------------------------------------------------------
				Info.Click.LandingID = LandingTemplateID
				Info.Click.PrelandingID = 0
				Info.Click.LocationLP = LandingTemplate

				defer Info.Click.Save()
				if cookieError != nil {
					defer c.SetCookie(utils.SaveCookieToUser(Info.Click.Hash, Info.Click.LocationLP))
				}

				// ----------------------------------------------------------------------------------------------------
				// STATS
				// ----------------------------------------------------------------------------------------------------
				config.TDSStatistic.RedirectRequest++ // add counter tick
				config.TDSStatistic.ProcessingTime += time.Since(start)

				if config.Cfg.Debug.Level > 1 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				if len(utils.ResponseAverage) < minimumStatCount {
					utils.ResponseAverage = append(utils.ResponseAverage, time.Since(start))
				} else {
					//defer runtime.GC()
					utils.ResponseAverage = utils.ResponseAverageDefault
				}
				// ----------------------------------------------------------------------------------------------------
				// FINAL
				// ----------------------------------------------------------------------------------------------------
				if !config.Cfg.Debug.Test {
				//	defer runtime.GC()
					return c.Redirect(302, LandingTemplate)
				} else {
					return c.Blob(200, "image/png", pixel)
				}
			}
		} else {
			config.TDSStatistic.IncorrectRequest++ // add counter tick
			// если нет клика или потока, то все привет
			msg := []byte(`{"code":400, "message":"Insuficient parameters supplied"}`)
			//defer runtime.GC()
			return c.JSONBlob(400, msg)
		}
	} else {
		// если нет редиски, то все привет
		msg := []byte(`{"code":500, "message":"No connection to RedisDB"}`)
		//defer runtime.GC()
		return c.JSONBlob(400, msg)
	}
}
