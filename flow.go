package main

import (
	"fmt"
	"github.com/labstack/echo"
	"math/rand"
	"metatds/config"
	"metatds/models"
	"metatds/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func clickBuild(c echo.Context) error {
	var Click models.ClickData

	if config.IsRedisAlive { // собираем данные для сейва в базу
		resultMap := utils.URIByMap(c, keyMap) // вот в этот массив

		PrelandID := strings.Join(resultMap["preland_id"], "")
		LandID := strings.Join(resultMap["land_id"], "")

		Click.Hash = strings.Join(resultMap["click_hash"], "")
		Click.FlowHash = strings.Join(resultMap["flow_hash"], "")

		resultMap["click_id"] = append(resultMap["click_id"], Click.Hash)

		if Click.Hash != "" && Click.FlowHash != "" {

			FlowID, _ := config.Redisdb.Get(Click.FlowHash + ":ID").Result()
			convertedID, _ := strconv.Atoi(FlowID)
			Click.FlowID = convertedID

			FlowWebMasterID, _ := config.Redisdb.Get(Click.FlowHash + ":WebMasterID").Result()
			convertedID, _ = strconv.Atoi(FlowWebMasterID)
			Click.WebMasterID = convertedID

			FlowWebMasterCurrencyID, _ := config.Redisdb.Get(Click.FlowHash + ":WebMasterCurrencyID").Result()
			convertedID, _ = strconv.Atoi(FlowWebMasterCurrencyID)
			Click.WebMasterCurrencyID = convertedID

			FlowOfferID, _ := config.Redisdb.Get(Click.FlowHash + ":OfferID").Result()
			convertedID, _ = strconv.Atoi(FlowOfferID)
			Click.OfferID = convertedID

			t := time.Now()

			Click.UserAgent = c.Request().UserAgent()
			Click.Time = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
				t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
			Click.URL = "http://" + config.Cfg.General.Host + c.Request().RequestURI
			Click.IP = c.Request().RemoteAddr
			Click.Referer = c.Request().Referer()

			Click.Sub1 = strings.Join(resultMap["sub1"], "")
			Click.Sub2 = strings.Join(resultMap["sub2"], "")
			Click.Sub3 = strings.Join(resultMap["sub3"], "")
			Click.Sub4 = strings.Join(resultMap["sub4"], "")
			Click.Sub5 = strings.Join(resultMap["sub5"], "")

			Prelands, _ := config.Redisdb.Keys(Click.FlowHash + ":preland:*").Result()
			PrelandingTemplate := ""

			for _, key := range Prelands {
				PrelandingTemplateID, _ := config.Redisdb.HGet(key, "id").Result()

				if PrelandingTemplateID == PrelandID {
					PrelandingTemplate, _ = config.Redisdb.HGet(key, "url").Result()

					for _, item := range keyMap {
						PrelandingTemplate = strings.Replace(PrelandingTemplate, fmt.Sprintf("{%s}", item),
							strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)

						fmt.Println("[ REPLACE ] = ", PrelandingTemplate)
					}
				}
			}

			fmt.Println("[ XXX ] = ", PrelandingTemplate)

			Click.LocationPL = PrelandingTemplate
			convertedID, _ = strconv.Atoi(PrelandID)
			Click.PrelandingID = convertedID
			Click.IsVisitedPL = 1

			Lands, _ := config.Redisdb.Keys(Click.FlowHash + ":land:*").Result()
			LandingTemplate := ""

			for _, key := range Lands {
				LandingTemplateID, _ := config.Redisdb.HGet(key, "id").Result()

				if LandingTemplateID == LandID {
					LandingTemplate, _ = config.Redisdb.HGet(key, "url").Result()

					for _, item := range keyMap {

						LandingTemplate = strings.Replace(LandingTemplate, fmt.Sprintf("{%s}", item),
							strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)

						fmt.Println("[ REPLACE ] = ", LandingTemplate)

					}
				}
			}

			Click.LocationLP = LandingTemplate
			convertedID, _ = strconv.Atoi(LandID)
			Click.LandingID = convertedID
			Click.IsVisitedLP = 1

			defer Click.Save()
		} else {
			msg := []byte(`{"code":400, "message":"No flow or click hashes found"}`)
			return c.JSONBlob(400, msg)
		}
	} else {
		// если нет редиски, то все привет
		msg := []byte(`{"code":400, "message":"No connection to RedisDB"}`)
		return c.JSONBlob(400, msg)
	}

	s := utils.JSONPretty(Click)

	return c.String(200, s)
}

/*
*
* Get Flow information in JSON ir Redirect to some selected location
* and Redirect handler all in one
*
*
 */

func flowHandler(c echo.Context) error {

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
		} else {
			Info.Click.Hash = Info.Click.GenerateCID()
		}

		// генерим СИД
		resultMap["click_hash"] = append(resultMap["click_hash"], Info.Click.Hash) // запишем сразу в наш массив
		resultMap["click_id"] = append(resultMap["click_id"], Info.Click.Hash)     // support for old version TDS

		Info.Flow = Info.Flow.GetInfo(strings.Join(resultMap["flow_hash"], "")) // получить всю инфу о потоке

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
				config.TDSStatistic.WorkTime += time.Since(start)

				if config.Cfg.Debug.Level > 0 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				if !config.Cfg.Debug.Test {
					return c.Redirect(302, LandingTemplate)
				} else {
					defer utils.WriteTestStatToFile(Info.Flow.Hash, Info.Click.Hash, LandingTemplate)

					return c.Blob(200, "image/png", pixel)
					//return c.String(200, "ok")
				}
			}
			// ----------------------------------------------------------------------------------------------------
			// PRELAND
			// редирект на пре-лендинг
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
				config.TDSStatistic.WorkTime += time.Since(start)

				if config.Cfg.Debug.Level > 0 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				if !config.Cfg.Debug.Test {
					return c.Redirect(302, PrelandingTemplate)
				} else {
					defer utils.WriteTestStatToFile(Info.Flow.Hash, Info.Click.Hash, LandingTemplate)

					return c.Blob(200, "image/png", pixel)
					//return c.String(200, "ok")
				}
			}
			// ----------------------------------------------------------------------------------------------------
			// JSON FORMAT
			// отдать данные потока в джейсоне красиво
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
				Info.Click.PrelandingID = PrelandingTemplateID
				Info.Click.LocationLP = LandingTemplate
				Info.Click.IsVisitedLP = 0
				Info.Click.LocationPL = PrelandingTemplate
				Info.Click.IsVisitedPL = 1
				//
				// TODO тут еще вопрос че показывать ленд или преленд
				//
				cookie.Path = Info.Click.LocationLP
				c.SetCookie(cookie)

				s := utils.JSONPretty(Info)

				defer Info.Click.Save()

				if config.Cfg.Debug.Level > 0 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				config.TDSStatistic.FlowInfoRequest++ // add counter tick
				config.TDSStatistic.WorkTime += time.Since(start)

				if !config.Cfg.Debug.Test {
					defer utils.WriteTestStatToFile(Info.Flow.Hash, Info.Click.Hash, LandingTemplate)
				}

				return c.String(200, s)

			} else {
				// OLD School
				// ----------------------------------------------------------------------------------------------------
				// если никаких ключей нет, то пробрасываем дальше (по старой схеме)
				// на первый выбраный домен из списка если несколько или на первый

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
				config.TDSStatistic.WorkTime += time.Since(start)

				if config.Cfg.Debug.Level > 0 {
					utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
				}

				if !config.Cfg.Debug.Test {
					return c.Redirect(302, LandingTemplate)
				} else {
					defer utils.WriteTestStatToFile(Info.Flow.Hash, Info.Click.Hash, LandingTemplate)

					return c.Blob(200, "image/png", pixel)
					//return c.String(200, "ok")
				}
			}
		} else {
			config.TDSStatistic.IncorrectRequest++ // add counter tick
			// если нет клика или потока, то все привет
			msg := []byte(`{"code":400, "message":"Insuficient parameters supplied"}`)
			return c.JSONBlob(400, msg)
		}
	} else {
		// если нет редиски, то все привет
		msg := []byte(`{"code":500, "message":"No connection to RedisDB"}`)
		return c.JSONBlob(400, msg)
	}
}
