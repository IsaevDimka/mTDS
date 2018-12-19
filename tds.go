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
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"metatds/config"
	"metatds/models"
)

const tdsModuleName = "tds.go"

// Карта ключей которые мы хотим получить
var keyMap = []string{"flow_hash", "click_hash", "sub1", "sub2", "sub3", "sub4", "sub5", "format",
	"click_id", "flow_id", "preland_id", "land_id"} // support for old version of TDS

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
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4/", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4/:sub5", clickBuild)
	router.GET("/c/build/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4/:sub5/", clickBuild)

	router.GET("/c/info/:click_hash", clickHandler)
	router.GET("/c/list", ListClickHandler)

	customServer := &http.Server{
		Addr:         ":" + strconv.Itoa(config.Cfg.General.Port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	customServer.SetKeepAlivesEnabled(false)

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
