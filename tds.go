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
	"fmt"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"math/rand"
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

// Карта ключей которые мы хотим получить
var keyMap = []string{"flow_hash", "click_hash", "sub1", "sub2", "sub3", "sub4", "sub5", "format",
	"click_id", "flow_id"} // support for old version of TDS

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
	router.GET("/c/:click_hash", clickHandler)
	router.GET("/c/all", allClickHandler)
	router.GET("/c/list", ListClickHandler)

	customServer := &http.Server{
		Addr:         ":" + strconv.Itoa(config.Cfg.General.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	router.HideBanner = true
	router.Logger.SetLevel(log.OFF)


	// run router
	if config.Cfg.General.Port != 0 {
//		router.Logger.Fatal(router.Start(":" + strconv.Itoa(config.Cfg.General.Port)))
		router.Logger.Fatal(router.StartServer(customServer))
	} else {
		//exit if not
		panic("[ERROR] Failed to obtain server port from settings.ini")
	}
}

/****************************************************************************************************
*
* Handlers for all responses
*
* allClicksHandler - returns JSON for all current Redis clicks for all Flows http://tds/c/all
* clickHandler - returns JSON for seleted Click HASH as param http://tds/c/PARAM
* flowHandler - 1. Finds Flow by Flow HASH as param
*               2. Generates CID
*               3. returns JSON for selected || redirect (302) depends on GET param "format"
*
*                http://tds/r/FLOW_HASH/sub1...sub2/other?get_params=other
*                http://tds/r/?flow_hash=FLOW_HASH&sub1=param...&sub5=param&other?get_params=other
*
/*****************************************************************************************************

/*
*
*
* Get all Clicks at current time from Redis
*
*
*/

func allClickHandler(c echo.Context) error {
	if config.IsRedisAlive {

		start := time.Now()

		var Click models.ClickData
		Clicks := make(map[string][]models.ClickData)
		FlowKeys, _ := config.Redisdb.Keys("*:Hash").Result()

		for _, item := range FlowKeys {
			CurrentFlowHash, _ := config.Redisdb.Get(item).Result()
			allClicks, _ := config.Redisdb.Keys(CurrentFlowHash + ":click:*").Result()

			for _, jtem := range allClicks {
				CurrentClickHash, _ := config.Redisdb.HGet(jtem, "Hash").Result()
				if config.Cfg.Debug.Level > 0 {
					utils.PrintInfo("Click item", CurrentFlowHash+" / "+CurrentClickHash, tdsModuleName)
				}
				Click = Click.GetInfo(CurrentClickHash)
				Clicks[CurrentFlowHash] = append(Clicks[CurrentFlowHash], Click)
			}
		}
		s := utils.JSONPretty(Clicks)

		//
		// TODO: Think about removing
		//

		if config.Cfg.Debug.Level > 0 {
			counter := 0
			for _, item := range Clicks {
				counter += len(item)
			}
			utils.PrintDebug("Count", strconv.Itoa(counter), tdsModuleName)
		}

		config.TDSStatistic.RedisStatRequest++ // add counter tick
		config.TDSStatistic.WorkTime += time.Since(start)

		utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)

		// TODO Здесь должен возвращаться контент тайп джейсон

		return c.String(200, s)
	} else {
		config.TDSStatistic.IncorrectRequest++ // add counter tick
		// если нет редиски, то все привет
		msg := []byte(`{"code":500, "message":"No connection to RedisDB"}`)
		return c.JSONBlob(400, msg)
	}
}

func ListClickHandler(c echo.Context) error {
	if config.IsRedisAlive {
		start := time.Now()

		var Clicks []models.ClickData
		var Click models.ClickData

		Keys, _ := config.Redisdb.Keys("*:click:*").Result()

		for _, item := range Keys {
			CurrentClickHash, _ := config.Redisdb.HGet(item, "Hash").Result()
			Click = Click.GetInfo(CurrentClickHash)
			Clicks = append(Clicks, Click)
		}

		s := utils.JSONPretty(Clicks)

		//
		// TODO: Think about removing
		//

		if config.Cfg.Debug.Level > 0 {
			counter := len(Clicks)
			utils.PrintDebug("Count", strconv.Itoa(counter), tdsModuleName)
		}

		config.TDSStatistic.RedisStatRequest++ // add counter tick
		config.TDSStatistic.WorkTime += time.Since(start)

		utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)

		// TODO Здесь должен возвращаться контент тайп джейсон

		return c.String(200, s)
	} else {
		config.TDSStatistic.IncorrectRequest++ // add counter tick
		// если нет редиски, то все привет
		msg := []byte(`{"code":500, "message":"No connection to RedisDB"}`)
		return c.JSONBlob(400, msg)
	}
}

/*
*
*
* Get Single Click from Redis
*
*
 */

func clickHandler(c echo.Context) error {
	if config.IsRedisAlive {
		var Click models.ClickData
		start := time.Now()

		resultMap := utils.URIByMap(c, keyMap)
		resultMap["click_hash"] = append(resultMap["click_hash"], Click.Hash) // запишем сразу в наш массив

		resultMap["click_id"] = append(resultMap["click_id"], Click.Hash) // support for old version TDS

		Click = Click.GetInfo(strings.Join(resultMap["click_hash"], ""))

		if Click != (models.ClickData{}) && config.Cfg.Debug.Level > 0 {
			utils.PrintDebug("Click info", Click, tdsModuleName)
		}

		s := utils.JSONPretty(Click)

		config.TDSStatistic.ClickInfoRequest++ // add counter tick
		config.TDSStatistic.WorkTime += time.Since(start)

		return c.String(200, s)
	} else {
		config.TDSStatistic.IncorrectRequest++ // add counter tick
		// если нет редиски, то все привет
		msg := []byte(`{"code":500, "message":"No connection to RedisDB"}`)
		return c.JSONBlob(400, msg)
	}
}

/*
*
* Get Flow information in JSON ir Redirect to some selected location
* and Redirect handler all in one
*
*
 */

func flowHandler(c echo.Context) error {
	if config.IsRedisAlive {
		start := time.Now()
		var Info InfoData
		var LandingTemplate, PrelandingTemplate string
		var LandingTemplateID, PrelandingTemplateID int
		resultMap := utils.URIByMap(c, keyMap) // вот в этот массив

		Info.Click.Hash = Info.Click.GenerateCID()

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
				Info.Click.Location = LandingTemplate

				cookie.Path = Info.Click.Location
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
				Info.Click.PrelandingID = PrelandingTemplateID
				Info.Click.Location = PrelandingTemplate

				cookie.Path = Info.Click.Location
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
				Info.Click.Location = LandingTemplate
				//
				// TODO тут еще вопрос че показывать ленд или преленд
				//
				cookie.Path = Info.Click.Location
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
				Info.Click.Location = LandingTemplate

				cookie.Path = Info.Click.Location
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
