package main

import (
	"fmt"
	"github.com/labstack/echo"
	"metatds/config"
	"metatds/models"
	"metatds/utils"
	"runtime"
	"strconv"
	"strings"
	"time"
)

/****************************************************************************************************
*
* Handlers for all responses
*
* allClicksHandler - returns JSON for all current Redis clicks for all Flows http://tds/c/all
* clickHandler - returns JSON for seleted Click HASH as param http://tds/c/PARAM
*
/*****************************************************************************************************/

//
// Get all Clicks at current time from Redis
//
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
		config.TDSStatistic.ProcessingTime += time.Since(start)

		utils.PrintInfo("Action elapsed time", time.Since(start), tdsModuleName)

		// TODO Здесь должен возвращаться контент тайп джейсон

		runtime.GC()
		return c.String(200, s)
	} else {
		config.TDSStatistic.IncorrectRequest++ // add counter tick
		// если нет редиски, то все привет
		msg := []byte(`{"code":500, "message":"No connection to RedisDB"}`)

		runtime.GC()
		return c.JSONBlob(400, msg)
	}
}

//
// Get Single Click from Redis
//
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
		config.TDSStatistic.ProcessingTime += time.Since(start)

		runtime.GC()
		return c.String(200, s)
	} else {
		config.TDSStatistic.IncorrectRequest++ // add counter tick
		// если нет редиски, то все привет
		msg := []byte(`{"code":500, "message":"No connection to RedisDB"}`)

		runtime.GC()
		return c.JSONBlob(400, msg)
	}
}

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
			runtime.GC()
		} else {
			msg := []byte(`{"code":400, "message":"No flow or click hashes found"}`)
			runtime.GC()
			return c.JSONBlob(400, msg)
		}
	} else {
		// если нет редиски, то все привет
		msg := []byte(`{"code":400, "message":"No connection to RedisDB"}`)
		return c.JSONBlob(400, msg)
	}

	s := utils.JSONPretty(Click)

	runtime.GC()
	return c.String(200, s)
}
