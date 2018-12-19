package main

import (
	"github.com/labstack/echo"
	"metatds/config"
	"metatds/models"
	"metatds/utils"
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
