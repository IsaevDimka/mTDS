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
	"bytes"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sevenNt/echo-pprof"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
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
var keyMap = []string{"flow_hash", "click_hash", "format", "f",
	"click_id", "flow_id", // support for old version of TDS
	"prelanding_id", "landing_id", "sub1", "sub2", "sub3", "sub4", "sub5",
	"utm_source", "utm_campaign", "utm_medium", "utm_content", "utm_term"} // support for old version of TDS

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

	//router.AutoTLSManager.HostPolicy = autocert.HostWhitelist("<DOMAIN>")
	// Cache certificates
	router.AutoTLSManager.Cache = autocert.DirCache("./.cache")

	// Middleware
	//router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	//avoid chrome to request favicon
	router.GET("/favicon.ico", func(c echo.Context) error {
		return c.Blob(200, "image/png", pixel)
		//return c.String(404, "not found") //nothing
	})

	// Routes
	router.GET("/r", flowHandler)
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
	router.GET("/c", clickBuild)
	router.GET("/c/", clickBuild)
	router.GET("/c/:flow_hash", clickBuild)
	router.GET("/c/:flow_hash/", clickBuild)
	router.GET("/c/:flow_hash/:click_hash", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id/", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id/:sub1", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4/", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4/:sub5", clickBuild)
	router.GET("/c/:flow_hash/:click_hash/:land_id/:preland_id/:sub1/:sub2/:sub3/:sub4/:sub5/", clickBuild)

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

	router.GET("/info/:click_hash", clickHandler)

	router.GET("/stat", GetSystemStatHandler)
	router.GET("/extstat", GetSystemExtendedStatHandler)
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

	// router.AutoTLSManager
	if config.Cfg.General.SSL {
		go func() {
			router.Logger.Fatal(router.StartTLS(":443", config.Cfg.General.SSLCert, config.Cfg.General.SSLKey))
		}()
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

func GetSystemStatHandler(c echo.Context) error {
	agent := c.Request().UserAgent()
	// защита от долбоебов
	if agent == "MetaDevAgent" {
		text := config.GetSystemStatistics()
		return c.HTML(200,
			"<html><head><title>TDS System statistics</title><script>"+
				"setInterval(function(){window.location.reload(true)},10000);"+
				"</script></head><body><pre>"+
				text+
				"</pre></body></html>")

	} else {
		return c.String(404, "Not found on server")
	}
}

func GetSystemExtendedStatHandler(c echo.Context) error {
	var dataForGraph string
	var dataKeys []int

	// Getting current stats 1 column
	text := config.GetSystemStatistics()

	// Reading file with template for graphic
	w := bytes.NewBuffer(nil)
	file, _ := os.Open("tmpl/sysstat.tmpl")
	_, _ = io.Copy(w, file)
	_ = file.Close()

	// Getting current stats from Redis saved before
	dataFromRedis, _ := config.Redisdb.HGetAll("SystemStatistic").Result()

	for i, _ := range dataFromRedis {
		convertedID, _ := strconv.Atoi(i)
		dataKeys = append(dataKeys, convertedID)
	}

	// needs to be sorted casue of HGetALL output issue
	sort.Ints(dataKeys)

	for _, item := range dataKeys {
		convertedID := strconv.Itoa(item)
		dataForGraph += dataFromRedis[convertedID] + ",\n"
	}

	// Just encode to json and print
	monitor := utils.MemMonitor()
	memoryStats := utils.JSONPretty(monitor)

	// getting template from var
	template := w.String()
	result := strings.Replace(template, "{{MEM}}", memoryStats, -1)
	result = strings.Replace(result, "{{DATA}}", dataForGraph, -1)
	result = strings.Replace(result, "{{SYSSTAT}}", text, -1)

	// returning resulting page
	return c.HTML(200, result)
}
