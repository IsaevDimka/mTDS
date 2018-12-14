/****************************************************************************************************
*
* Main TDS module, special for Meta CPA, Ltd.
* by Michael S. Merzlyakov AFKA predator_pc@02122018
* version v2.0.3
*
* created at 02122018
* last edit: 13122018
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
)

const tdsModuleName = "tds.go"

// Карта ключей которые мы хотим получить
var keyMap = []string{"flow_hash", "click_hash", "sub1", "sub2", "sub3", "sub4", "sub5", "format"}

// Тип для
type InfoData struct {
	Flow  FlowData  // модель потока
	Click ClickData // модель клика
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
	if cfg.General.Port != 0 {
		router.Logger.Fatal(router.Start(":" + strconv.Itoa(cfg.General.Port)))
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
	var Click ClickData
	Clicks := make(map[string][]ClickData)
	FlowKeys, _ := redisdb.Keys("*:FlowHash").Result()

	for _, item := range FlowKeys {
		CurrentFlowHash, _ := redisdb.Get(item).Result()
		allClicks, _ := redisdb.Keys(CurrentFlowHash + ":click:*").Result()

		for _, jtem := range allClicks {
			CurrentClickHash, _ := redisdb.HGet(jtem, "Hash").Result()
			if cfg.Debug.Level > 0 {
				printInfo("Click item", CurrentFlowHash+" / "+CurrentClickHash, tdsModuleName)
			}
			Click = Click.getInfo(CurrentClickHash)
			Clicks[CurrentFlowHash] = append(Clicks[CurrentFlowHash], Click)
		}
	}
	s := JSONPretty(Clicks)

	//
	// TODO: Think about removing
	//

	if cfg.Debug.Level > 0 {
		counter := 0
		for _, item := range Clicks {
			counter += len(item)
		}
		printDebug("Count", strconv.Itoa(counter), tdsModuleName)
	}
	return c.String(200, s)
}

/*
* Get Single Click from Redis
 */

func clickHandler(c echo.Context) error {
	var Click ClickData

	resultMap := URIByMap(c, keyMap)
	resultMap["click_hash"] = append(resultMap["click_hash"], Click.Hash) // запишем сразу в наш массив
	Click = Click.getInfo(strings.Join(resultMap["click_hash"], ""))

	if Click != (ClickData{}) && cfg.Debug.Level > 0 {
		printDebug("Click info", Click, tdsModuleName)
	}

	s := JSONPretty(Click)
	return c.String(200, s)
}

/*
* Get Flow information in JSON ir Redirect to some selected location
* and Redirect handler all in one
 */

func flowHandler(c echo.Context) error {
	var Info InfoData
	var LandingTemplate, PrelandingTemplate, LandingTemplateID, PrelandingTemplateID string
	resultMap := URIByMap(c, keyMap) // вот в этот массив

	Info.Click.Hash = Info.Click.generateCID()                                 // генерим СИД
	resultMap["click_hash"] = append(resultMap["click_hash"], Info.Click.Hash) // запишем сразу в наш массив
	Info.Flow = Info.Flow.getInfo(strings.Join(resultMap["flow_hash"], ""))    // получить всю инфу о потоке

	t := time.Now() // собираем данные для сейва в базу

	Info.Click.UserAgent = c.Request().UserAgent()
	Info.Click.Time = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	Info.Click.URL = "http://" + cfg.General.Host + c.Request().RequestURI

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
			LandingTemplateID = Info.Flow.Lands[Random].ID
		} else {
			msg := []byte(`{"code":400, "message":"No landing templates found"}`)
			return c.JSONBlob(400, msg)
		}

		if len(Info.Flow.Prelands) > 0 {
			Random := rand.Intn(len(Info.Flow.Prelands))
			PrelandingTemplate = Info.Flow.Prelands[Random].URL // получаем рандомный урл преленда
			PrelandingTemplateID = Info.Flow.Prelands[Random].ID
		}

		// Если дебаг то печатаем все это добро
		if cfg.Debug.Level > 1 && len(Info.Flow.Prelands) > 0 && len(Info.Flow.Lands) > 0 {
			printDebug("Land", resultMap, tdsModuleName)
			printDebug("Land", Info.Flow.Lands[rand.Intn(len(Info.Flow.Lands))].URL, tdsModuleName)
			printDebug("Preland", Info.Flow.Prelands[rand.Intn(len(Info.Flow.Prelands))].URL, tdsModuleName)
		}

		// выбираем по формату, что будем отдавать
		// редирект на лендинг
		if strings.Join(resultMap["format"], "") == "land" {
			for _, item := range keyMap {
				LandingTemplate = strings.Replace(LandingTemplate, fmt.Sprintf("{%s}", item),
					strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
			}

			if cfg.Debug.Level > 1 {
				printInfo("Result", LandingTemplate, tdsModuleName)
			}

			convertedID, _ := strconv.Atoi(LandingTemplateID)
			Info.Click.LandingID = convertedID
			Info.Click.PrelandingID = 0
			Info.Click.Location = LandingTemplate

			defer Info.Click.save()
			return c.Redirect(302, LandingTemplate)
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

			if cfg.Debug.Level > 1 {
				printInfo("Result", PrelandingTemplate, tdsModuleName)
			}

			convertedID, _ := strconv.Atoi(PrelandingTemplateID)
			Info.Click.LandingID = 0
			Info.Click.PrelandingID = convertedID
			Info.Click.Location = PrelandingTemplate

			defer Info.Click.save()
			return c.Redirect(302, PrelandingTemplate)
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

			s := JSONPretty(Info)

			defer Info.Click.save()
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

			defer Info.Click.save()
			return c.Redirect(302, LandingTemplate)
		}
	} else {
		// если нет клика или потока, то все привет
		msg := []byte(`{"code":400, "message":"Insuficient parameters supplied"}`)
		return c.JSONBlob(400, msg)
	}

	// должны возвращать по формату функции, на самом деле никогда не выполнится
	msg := []byte(`{"code":200, "message":"Destination unreacheable"}`)
	return c.JSONBlob(200, msg)
}

/*
* Just for debug
 */

func testHandler(c echo.Context) error {
	var Click ClickData
	t := time.Now()

	Click.UserAgent = c.Request().UserAgent()
	Click.Time = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	Click.IP = c.Request().RemoteAddr
	Click.URL = "http://" + cfg.General.Host + c.Request().RequestURI
	Click.Referer = c.Request().Referer()

	fmt.Println("[CLICK]", Click)

	return c.String(200, "ok")
}
