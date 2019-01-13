/****************************************************************************************************
*
* Handlers for cliks
*
* allClicksHandler - returns JSON for all current Redis clicks for all Flows http://tds/c/all
* clickHandler - returns JSON for seleted Click HASH as param http://tds/c/PARAM
* special for Meta CPA, Ltd.
* by Michael S. Merzlyakov AFKA predator_pc@10012019
* version v2.0.3
*
* created at 04122018
* last edit: 13012018
*
*****************************************************************************************************/

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
		return c.String(200, s)
	} else {
		config.TDSStatistic.IncorrectRequest++ // add counter tick
		// если нет редиски, то все привет
		msg := []byte(`{"code":500, "message":"No connection to RedisDB"}`)
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
	var Flow models.FlowData

	if config.IsRedisAlive { // собираем данные для сейва в базу
		resultMap := utils.URIByMap(c, keyMap) // вот в этот массив

		PrelandID := strings.Join(resultMap["preland_id"], "")
		LandID := strings.Join(resultMap["land_id"], "")

		Click.Hash = strings.Join(resultMap["click_hash"], "")
		Click.FlowHash = strings.Join(resultMap["flow_hash"], "")

		// Костыли для старого стиля обращений
		if Click.Hash == "" {
			Click.Hash = strings.Join(resultMap["click_id"], "")
		}
		if Click.FlowHash == "" {
			Click.FlowHash = strings.Join(resultMap["flow_id"], "")
		}

		resultMap["click_id"] = append(resultMap["click_id"], Click.Hash)

		if Click.Hash != "" && Click.FlowHash != "" {

			Flow = Flow.GetInfo(Click.FlowHash) // получить всю инфу о потоке
			Click.FlowID = Flow.ID
			Click.WebMasterID = Flow.WebMasterID
			Click.WebMasterCurrencyID = Flow.WebMasterCurrencyID
			Click.OfferID = Flow.OfferID

			Click.UserAgent = c.Request().UserAgent()
			Click.Time = utils.CURRENT_TIMESTAMP
			Click.URL = "http://" + config.Cfg.General.Host + c.Request().RequestURI
			Click.IP = c.Request().RemoteAddr
			Click.Referer = c.Request().Referer()

			Click.Sub1 = strings.Join(resultMap["sub1"], "")
			Click.Sub2 = strings.Join(resultMap["sub2"], "")
			Click.Sub3 = strings.Join(resultMap["sub3"], "")
			Click.Sub4 = strings.Join(resultMap["sub4"], "")
			Click.Sub5 = strings.Join(resultMap["sub5"], "")

			// TODO: Возможна оптимизация по прямому обращению к индексу искомого ленда или преленда
			Prelands, _ := config.Redisdb.HGetAll(Click.FlowHash + ":prelands").Result()
			PrelandingTemplate := ""

			for PrelandingTemplateID, key := range Prelands {
				if PrelandingTemplateID == PrelandID {
					PrelandingTemplate = key
					for _, item := range keyMap {
						PrelandingTemplate = strings.Replace(PrelandingTemplate, fmt.Sprintf("{%s}", item),
							strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
					}
				}
			}

			Click.LocationPL = PrelandingTemplate
			convertedID, _ := strconv.Atoi(PrelandID)
			Click.PrelandingID = convertedID
			Click.IsVisitedPL = 1

			// TODO: Возможна оптимизация по прямому обращению к индексу искомого ленда или преленда
			Lands, _ := config.Redisdb.HGetAll(Click.FlowHash + ":lands").Result()
			LandingTemplate := ""

			for LandingTemplateID, key := range Lands {
				if LandingTemplateID == LandID {
					LandingTemplate = key
					for _, item := range keyMap {
						LandingTemplate = strings.Replace(LandingTemplate, fmt.Sprintf("{%s}", item),
							strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
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
