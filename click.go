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
	"bytes"
	"fmt"
	"github.com/labstack/echo"
	"github.com/predatorpc/durafmt"
	"io"
	"metatds/config"
	"metatds/models"
	"metatds/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//
// Get Single Click from Redis
//
func clickHandler(c echo.Context) error {

	if config.IsRedisAlive {

		var Click models.ClickData
		start := time.Now()

		resultMap, _ := utils.URIByMap(c, keyMap)
		resultMap["click_hash"] = append(resultMap["click_hash"], Click.Hash) // запишем сразу в наш массив
		resultMap["click_id"] = append(resultMap["click_id"], Click.Hash)     // support for old version TDS

		if strings.Join(resultMap["format"], "") == "local" || strings.Join(resultMap["f"], "") == "local" {
			Click = Click.GetInfo(strings.Join(resultMap["click_hash"], ""))

			if Click != (models.ClickData{}) && config.Cfg.Debug.Level > 0 {
				utils.PrintDebug("Click info", Click, tdsModuleName)
			}

			data := utils.JSONPretty(Click)

			config.TDSStatistic.ClickBuildRequest++ // add counter tick
			config.TDSStatistic.ProcessingTime += time.Since(start)

			return c.String(200, data+"\nTotal elapsed: "+
				durafmt.Parse(time.Since(start)).String(durafmt.DF_LONG))
		}
		if strings.Join(resultMap["format"], "") == "remote" || strings.Join(resultMap["f"], "") == "remote" {

			body := bytes.NewBuffer(nil)
			req, err := http.Get("http://116.202.27.130/hit/" + strings.Join(resultMap["click_hash"], ""))

			if err != nil {
				recover()
				fmt.Println("[ ERROR] can't create request 0: ", err)
			}

			if req != nil {
				// setting header in case of API request is not NIL
				req.Header.Set("Connection", "close")
				// reading the body
				_, _ = io.Copy(body, req.Body)
				// closing anyway now
				// defer is not needed cause we get an exception before
				_ = req.Body.Close()
			}

			config.TDSStatistic.ClickBuildRequest++ // add counter tick
			config.TDSStatistic.ProcessingTime += time.Since(start)

			return c.String(200, body.String()+"\nTotal elapsed: "+
				durafmt.Parse(time.Since(start)).String(durafmt.DF_LONG))
		} else {
			config.TDSStatistic.IncorrectRequest++ // add counter tick
			// если нет редиски, то все привет
			msg := []byte(`{"code":400, "message":"Please specify format local/remote"}`)
			return c.JSONBlob(400, msg)
		}
	} else {
		config.TDSStatistic.IncorrectRequest++ // add counter tick
		// если нет редиски, то все привет
		msg := []byte(`{"code":400, "message":"No connection to RedisDB"}`)
		return c.JSONBlob(400, msg)
	}
}

//
// Resulting click when user goes on landing or it could be a post-back without answer
// in case if we want to tell that user achieved a goal
//
func clickBuild(c echo.Context) error {
	var Click models.ClickData
	var Flow models.FlowData

	resultMap, _ := utils.URIByMap(c, keyMap) // вот в этот массив
	// если редис жив
	if config.IsRedisAlive { // собираем данные для сейва в базу

		// айдишники лендов и прелендов для построения образа клика
		PrelandID := strings.Join(resultMap["prelanding_id"], "")
		LandID := strings.Join(resultMap["landing_id"], "")

		// Костыли для старого стиля обращений TDS v1.0
		Click.Hash = strings.Join(resultMap["click_hash"], "")
		if Click.Hash == "" {
			Click.Hash = strings.Join(resultMap["click_id"], "")
		} else {
			resultMap["click_id"] = append(resultMap["click_id"], Click.Hash)
		}

		Click.FlowHash = strings.Join(resultMap["flow_hash"], "")
		// Костыли для старого стиля обращений TDS v1.0
		if Click.FlowHash == "" {
			Click.FlowHash = strings.Join(resultMap["flow_id"], "")
		}

		if Click.Hash != "" && Click.FlowHash != "" {
			Flow = Flow.GetInfo(Click.FlowHash) // получить всю инфу о потоке
			Click.Time = utils.CURRENT_TIMESTAMP

			Click.FlowID = Flow.ID
			Click.WebMasterID = Flow.WebMasterID
			Click.WebMasterCurrencyID = Flow.WebMasterCurrencyID
			Click.OfferID = Flow.OfferID

			// Получаем проброшенные заголовки т.к. CURL нам просто так не отдаст клиентские
			XRealIP := c.Request().Header.Get("X-Real-IP")
			if XRealIP != "" {
				Click.IP = XRealIP
			} else {
				Click.IP = c.Request().RemoteAddr
			}

			Click.UserAgent = c.Request().UserAgent()
			Click.Referer = c.Request().Referer()

			Click.Sub1 = strings.Join(resultMap["sub1"], "")
			Click.Sub2 = strings.Join(resultMap["sub2"], "")
			Click.Sub3 = strings.Join(resultMap["sub3"], "")
			Click.Sub4 = strings.Join(resultMap["sub4"], "")
			Click.Sub5 = strings.Join(resultMap["sub5"], "")

			// Получаем все прелендинги
			Prelands, _ := config.Redisdb.HGetAll(Click.FlowHash + ":prelands").Result()
			// Берем тот который нам сказали
			PrelandingTemplate := Prelands[PrelandID]
			// Проверяем его наличие
			if PrelandingTemplate != "" {
				for _, item := range keyMap {
					PrelandingTemplate = strings.Replace(PrelandingTemplate, fmt.Sprintf("{%s}", item),
						strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
				}
			} else {
				// в случае ошибки говорим веб-мастеру попрравить
				msg := []byte(`{"code":400, "message":"No prelanding provided with such parameters"}`)
				return c.JSONBlob(400, msg)
			}

			// Устанавливаем данные в клик
			Click.LocationPL = PrelandingTemplate
			convertedID, _ := strconv.Atoi(PrelandID)
			Click.PrelandingID = convertedID
			// Признак посещения преленда
			Click.IsVisitedPL = 1

			// Получаем все лендинги
			Lands, _ := config.Redisdb.HGetAll(Click.FlowHash + ":lands").Result()
			// Берем тот который нам сказали
			LandingTemplate := Lands[LandID]
			// Проверяем что он вообще есть
			if LandingTemplate != "" {
				for _, item := range keyMap {
					LandingTemplate = strings.Replace(LandingTemplate, fmt.Sprintf("{%s}", item),
						strings.Trim(fmt.Sprintf("%s", resultMap[item]), " ]["), 1)
				}
			} else {
				// в случае ошибки говорим веб-мастеру попрравить
				msg := []byte(`{"code":400, "message":"No landing provided with such parameters"}`)
				return c.JSONBlob(400, msg)
			}

			// Устанавливаем данные в клик
			Click.LocationLP = LandingTemplate
			convertedID, _ = strconv.Atoi(LandID)
			Click.LandingID = convertedID
			// Признак посещения ленда
			Click.IsVisitedLP = 1

			config.TDSStatistic.ClickBuildRequest++

			// отправляем клик в редис
			defer Click.Save()
		} else {
			// если нам не передали клика или потока, то шлем лесом
			msg := []byte(`{"code":400, "message":"No flow or click hashes found"}`)
			return c.JSONBlob(400, msg)
		}
	} else {
		// если нет редиски, то все привет
		msg := []byte(`{"code":400, "message":"No connection to RedisDB"}`)
		return c.JSONBlob(400, msg)
	}

	// Finally return data or redirect in 2 case CURL may not receive answer
	if strings.Join(resultMap["format"], "") == "lp" || strings.Join(resultMap["f"], "") == "lp" {
		return c.Redirect(302, Click.LocationLP)
	} else {
		s := utils.JSONPretty(Click)
		return c.String(200, s)
	}
}
