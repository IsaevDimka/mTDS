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
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"metatds/config"
	"metatds/models"
	"metatds/utils"
)

const tdsModuleName = "tds.go"

// Карта ключей которые мы хотим получить
var keyMap = []string{"flow_hash", "click_hash", "sub1", "sub2", "sub3", "sub4", "sub5", "format",
	"click_id", "flow_id"} // support for old version of TDS

// Тип для
type InfoData struct {
	Flow  models.FlowData  // модель потока
	Click models.ClickData // модель клика
}

/*
* Main GO package handler
 */

func main() {
	// Echo instance
	router := echo.New()
	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	//avoid chrome to request favicon
	router.GET("/favicon.ico", func(c echo.Context) error {
		return c.String(404, "not found") //nothing
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
	router.GET("/test", testHandler)
	// run router
	if config.Cfg.General.Port != 0 {
		router.Logger.Fatal(router.Start(":" + strconv.Itoa(config.Cfg.General.Port)))
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
* testhandler - returns anything you described in it, just for tests
*
/****************************************************************************************************

/*
* Get all Clicks at current time from Redis
*/

func allClickHandler(c echo.Context) error {
	start := time.Now()

	var Click models.ClickData
	Clicks := make(map[string][]models.ClickData)
	FlowKeys, _ := config.Redisdb.Keys("*:FlowHash").Result()

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
	utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)

	// TODO Здесь должен возвращаться контент тайп джейсон

	return c.String(200, s)
}

/*
* Get Single Click from Redis
 */

func clickHandler(c echo.Context) error {
	var Click models.ClickData

	resultMap := utils.URIByMap(c, keyMap)
	resultMap["click_hash"] = append(resultMap["click_hash"], Click.Hash) // запишем сразу в наш массив

	resultMap["click_id"] = append(resultMap["click_id"], Click.Hash) // support for old version TDS

	Click = Click.GetInfo(strings.Join(resultMap["click_hash"], ""))

	if Click != (models.ClickData{}) && config.Cfg.Debug.Level > 0 {
		utils.PrintDebug("Click info", Click, tdsModuleName)
	}

	s := utils.JSONPretty(Click)

	config.TDSStatistic.ClickInfoRequest++ // add counter tick

	return c.String(200, s)
}

/*
* Get Flow information in JSON ir Redirect to some selected location
* and Redirect handler all in one
 */

func flowHandler(c echo.Context) error {
	start := time.Now()

	var Info InfoData
	var LandingTemplate, PrelandingTemplate, LandingTemplateID, PrelandingTemplateID string
	resultMap := utils.URIByMap(c, keyMap) // вот в этот массив

	Info.Click.Hash = Info.Click.GenerateCID()                                 // генерим СИД
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
			LandingTemplate = Info.Flow.Lands[Random].URL // получаем рандомный урл ленда
			LandingTemplateID = strconv.Itoa(Info.Flow.Lands[Random].ID)
		} else {
			msg := []byte(`{"code":400, "message":"No landing templates found"}`)
			return c.JSONBlob(400, msg)
		}

		if len(Info.Flow.Prelands) > 0 {
			Random := rand.Intn(len(Info.Flow.Prelands))
			PrelandingTemplate = Info.Flow.Prelands[Random].URL // получаем рандомный урл преленда
			PrelandingTemplateID = strconv.Itoa(Info.Flow.Prelands[Random].ID)
		}

		// Если дебаг то печатаем все это добро
		if config.Cfg.Debug.Level > 1 && len(Info.Flow.Prelands) > 0 && len(Info.Flow.Lands) > 0 {
			utils.PrintDebug("Land", resultMap, tdsModuleName)
			utils.PrintDebug("Land", Info.Flow.Lands[rand.Intn(len(Info.Flow.Lands))].URL, tdsModuleName)
			utils.PrintDebug("Preland", Info.Flow.Prelands[rand.Intn(len(Info.Flow.Prelands))].URL, tdsModuleName)
		}

		// выбираем по формату, что будем отдавать
		// редирект на лендинг
		if strings.Join(resultMap["format"], "") == "land" {
			for _, item := range keyMap {
				LandingTemplate = strings.Replace(LandingTemplate, fmt.Sprintf("{%s}", item),
					strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
			}

			if config.Cfg.Debug.Level > 1 {
				utils.PrintInfo("Result", LandingTemplate, tdsModuleName)
			}

			convertedID, _ := strconv.Atoi(LandingTemplateID)
			Info.Click.LandingID = convertedID
			Info.Click.PrelandingID = 0
			Info.Click.Location = LandingTemplate

			defer Info.Click.Save()

			config.TDSStatistic.RedirectRequest++ // add counter tick

			if config.Cfg.Debug.Level > 0 {
				utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
			}

			if !config.Cfg.Debug.Test {
				return c.Redirect(302, LandingTemplate)
			} else {
				return c.String(200, "ok")
			}
		}

		// редирект на пре-лендинг
		if strings.Join(resultMap["format"], "") == "preland" {

			//break
			if len(Info.Flow.Prelands) <= 0 {
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

			convertedID, _ := strconv.Atoi(PrelandingTemplateID)
			Info.Click.LandingID = 0
			Info.Click.PrelandingID = convertedID
			Info.Click.Location = PrelandingTemplate

			defer Info.Click.Save()
			config.TDSStatistic.RedirectRequest++ // add counter tick

			if config.Cfg.Debug.Level > 0 {
				utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
			}

			if !config.Cfg.Debug.Test {
				return c.Redirect(302, PrelandingTemplate)
			} else {
				return c.String(200, "ok")
			}
		}

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

			s := utils.JSONPretty(Info)

			defer Info.Click.Save()

			if config.Cfg.Debug.Level > 0 {
				config.TDSStatistic.FlowInfoRequest++ // add counter tick
			}

			utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)

			return c.String(200, s)

		} else {
			// если никаких ключей нет, то пробрасываем дальше (по старой схеме)
			// на первый выбраный домен из списка если несколько или на первый
			for _, item := range keyMap {
				LandingTemplate = strings.Replace(LandingTemplate, fmt.Sprintf("{%s}", item),
					strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
			}

			convertedID, _ := strconv.Atoi(LandingTemplateID)
			Info.Click.LandingID = 0
			Info.Click.PrelandingID = convertedID
			Info.Click.Location = LandingTemplate

			defer Info.Click.Save()
			config.TDSStatistic.RedirectRequest++ // add counter tick

			if config.Cfg.Debug.Level > 0 {
				utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)
			}

			if !config.Cfg.Debug.Test {
				return c.Redirect(302, LandingTemplate)
			} else {
				return c.String(200, "ok")
			}
		}
	} else {
		// если нет клика или потока, то все привет
		msg := []byte(`{"code":400, "message":"Insuficient parameters supplied"}`)
		return c.JSONBlob(400, msg)
	}

	// должны возвращать по формату функции, на самом деле никогда не выполнится
	// msg := []byte(`{"code":200, "message":"Destination unreacheable"}`)
	// return c.JSONBlob(200, msg)
}

/*
* Just for debug
 */

func testHandler(c echo.Context) error {
	var Click models.ClickData
	t := time.Now()

	Click.UserAgent = c.Request().UserAgent()
	Click.Time = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	Click.IP = c.Request().RemoteAddr
	Click.URL = "http://" + config.Cfg.General.Host + c.Request().RequestURI
	Click.Referer = c.Request().Referer()

	fmt.Println("[CLICK]", Click)

	return c.String(200, "ok")
}
