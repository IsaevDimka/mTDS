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
	"github.com/sevenNt/echo-pprof"
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"metatds/config"
	"metatds/models"
	"metatds/utils"
)

const tdsModuleName = "tds.go"
const timestampFile = "last.update.time"
const defaultStartOfEpoch = "946684800"

// Карта ключей которые мы хотим получить
var keyMap = []string{"flow_hash", "click_hash", "sub1", "sub2", "sub3", "sub4", "sub5", "format",
	"click_id", "flow_id", "preland_id", "land_id"} // support for old version of TDS

// пиксель для тестов
var pixel = []byte(`data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HgAGgwJ/lK3Q6wAAAABJRU5ErkJggg==`)

// Тип для полного потока с кликами
type InfoData struct {
	Flow  models.FlowData  // модель потока
	Click models.ClickData // модель клика
}

//
// Main GO package handler startup and settings
//
func main() {
	// let's update our flows on start
	UpdateFlowsListChan()

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

	router.GET("/stat", GetSystemStatkHandler)
	router.GET("/conf", GetSystemConfHandler)

	router.GET("/free", func(c echo.Context) error {
		debug.FreeOSMemory()
		runtime.GC()
		return c.String(200, "ok")
	})

	router.GET("/memstat", func(c echo.Context) error {

		m := utils.MemMonitor()
		// Just encode to json and print
		return c.String(200, utils.JSONPretty(m))
	})

	customServer := &http.Server{
		Addr:         ":" + strconv.Itoa(config.Cfg.General.Port),
		ReadTimeout:  time.Duration(1+config.Cfg.General.HTTPTimeout) * time.Second,
		WriteTimeout: time.Duration(1+config.Cfg.General.HTTPTimeout) * time.Second,
		IdleTimeout:  time.Duration(1+config.Cfg.General.HTTPTimeout) * time.Second,
	}

	customServer.SetKeepAlivesEnabled(false)
	router.HideBanner = true
	router.Logger.SetLevel(log.OFF)

	if config.Cfg.Debug.Level > 0 {
		echopprof.Wrap(router)
	}

	// run router
	if config.Cfg.General.Port != 0 {
		// regular server
		// router.Logger.Fatal(router.Start(":" + strconv.Itoa(config.Cfg.General.Port)))

		// custom server
		router.Logger.Fatal(router.StartServer(customServer))
	} else {
		//exit if not
		panic("[ERROR] Failed to obtain server port from settings.ini")
	}

}

func GetSystemConfHandler(c echo.Context) error {
	agent := c.Request().UserAgent()
	// защита от долбоебов
	if agent == "MetaDevAgent" {
		text := config.GetSystemConfiguration()
		return c.HTML(200,
			"<html><head><title>TDS System statistics</title><script>"+
				"setInterval(function(){window.location.reload(true)},5000);"+
				"</script></head><body><pre>"+
				text+
				"</pre></body></html>")

	} else {
		return c.String(404, "Not found on server")
	}
}

func GetSystemStatkHandler(c echo.Context) error {
	agent := c.Request().UserAgent()
	// защита от долбоебов
	if agent == "MetaDevAgent" {
		text := config.GetSystemStatistics()
		return c.HTML(200,
			"<html><head><title>TDS System statistics</title><script>"+
				"setInterval(function(){window.location.reload(true)},5000);"+
				"</script></head><body><pre>"+
				text+
				"</pre></body></html>")

	} else {
		return c.String(404, "Not found on server")
	}
}
