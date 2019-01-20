/****************************************************************************************************
*
* System statistics module / visualization, special for Meta CPA, Ltd.
* by Michael S. Merzlyakov AFKA predator_pc@09012019
* version v2.0.5
*
* created at 04122018
* last edit: 20012019
*
*****************************************************************************************************/

package config

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/predatorpc/durafmt"
	"io/ioutil"
	"math"
	"metatds/utils"
	"os"
	"runtime"
	"strconv"
	"time"
)

var RPSStat []int
var RPSStatDefault []int

const minimumStatCountRPS = 100
const sysstatModuleName = "sysstat.go"

func GetSystemStatistics() string {
	var text = "no stat"
	var memory runtime.MemStats
	var duration time.Duration // current duration & uptime
	var uptime, processingTime, memoryUsageGeneral, memoryUsagePrivate, avgReq string
	var openedFiles = "0"

	// Setting up counter if not exists
	StatisticCounter, err := Redisdb.Get("StatisticCounter").Result()
	if err != nil {
		_ = Redisdb.Set("StatisticCounter", 0, 0).Err()
		StatisticCounter = "0"
	}

	// Если статистики нет, мы должны инициализировать структуру с ней
	if TDSStatistic != (utils.TDSStats{}) {

		// AVERAGE RPS STATS  -------------------------------------------------------------------------------------------
		//
		currentRPSstart := TDSStatistic.RedirectRequest
		time.Sleep(1 * time.Second)
		currentRPS := TDSStatistic.RedirectRequest - currentRPSstart

		if len(RPSStat) < minimumStatCountRPS {
			RPSStat = append(RPSStat, currentRPS)
		} else {
			RPSStat = RPSStatDefault
		}
		averageRPS := RPSAverage(RPSStat)

		// PROCESSING AND UPTIME ----------------------------------------------------------------------------------------
		//
		duration = 60 * time.Minute
		if time.Since(UpTime) < duration {
			uptime = durafmt.Parse(time.Since(UpTime)).String(durafmt.DF_LONG)
			processingTime = durafmt.Parse(TDSStatistic.ProcessingTime).String(durafmt.DF_LONG)
		} else {
			uptime = durafmt.Parse(time.Since(UpTime)).String(durafmt.DF_MIDDLE)
			processingTime = durafmt.Parse(TDSStatistic.ProcessingTime).String(durafmt.DF_MIDDLE)
		}

		// MEMORY USAGE -------------------------------------------------------------------------------------------------
		//
		runtime.ReadMemStats(&memory)
		RealDetectedGeneral := memory.Sys + memory.HeapSys + memory.HeapAlloc + memory.HeapInuse - memory.Alloc
		RealDetectedPrivate := memory.HeapSys - memory.Alloc
		memoryUsageGeneral = strconv.FormatUint(utils.BToMb(RealDetectedGeneral), 10)
		memoryUsagePrivate = strconv.FormatUint(utils.BToMb(RealDetectedPrivate), 10)

		// GET OPENED SOCKET STATS ---------------------------------------------------------------------------------------
		// no need to do this for windows becasue of WSA working differently
		//
		if Cfg.General.OS == "linux" || Cfg.General.OS == "unix" {

			pid := strconv.Itoa(os.Getpid())
			fds, e := ioutil.ReadDir("/proc/" + pid + "/fd")

			if e != nil && Cfg.Debug.Level > 0 {
				utils.PrintError("Error", "reading process directory failed", sysstatModuleName)
			} else {
				if Cfg.Debug.Level > 1 {
					utils.PrintInfo("Reding stats PID", pid, initModuleName)
				}
			}

			if len(fds) > 0 {
				openedFiles = strconv.Itoa(len(fds))
			}
		}

		// RESPONSE AVERAGE ----------------------------------------------------------------------------------------------
		// Если среднее время ответа меньше чем миллисекунда так и напишем
		//
		dur := DurationAverage(utils.ResponseAverage)
		if dur < time.Duration(1*time.Millisecond) {
			avgReq = "< 1 msec"
		} else { // иначе напишем по человечески
			avgReq = durafmt.Parse(dur).String(durafmt.DF_LONG)
		}
		// ---------------------------------------------------------------------------------------------------------------

		// TODO: Тут надо думать как нам считать уники по первым запросам или нет
		// TODO: очевидно что куки реквесты мы удаляем из уников остается
		// TODO: Первыичный запрос и запрос JSON для потока, я думаю что первичный самый важный в итоге
		// uniqueRequests := (TDSStatistic.ClickBuildRequest + TDSStatistic.FlowInfoRequest + TDSStatistic.RedirectRequest) - TDSStatistic.CookieRequest
		uniqueRequests := TDSStatistic.RedirectRequest - TDSStatistic.CookieRequest

		//
		// auto update statistics in graph
		// and auto slide system stats removing first element when the count achieved `minimumStatCountRPS`
		//
		convertedID, _ := strconv.Atoi(StatisticCounter)
		if convertedID < minimumStatCountRPS {
			//			fmt.Println("Appending", convertedID)
			_ = Redisdb.HSet("SystemStatistic", StatisticCounter, "["+StatisticCounter+","+
				fmt.Sprintf("%.0f", (math.Round(dur.Seconds()*1000)))+","+
				strconv.Itoa(averageRPS)+","+strconv.Itoa(currentRPS)+"]").Err()
		} else {
			convertedID = convertedID - minimumStatCountRPS
			for i := 0; i < convertedID; i++ {
				_ = Redisdb.HDel("SystemStatistic", strconv.Itoa(i)).Err()
			}
			_ = Redisdb.HSet("SystemStatistic", StatisticCounter, "["+StatisticCounter+","+
				fmt.Sprintf("%.0f", (math.Round(dur.Seconds()*1000)))+","+
				strconv.Itoa(averageRPS)+","+strconv.Itoa(currentRPS)+"]").Err()
		}

		// allow to overwrite statistics
		text = "\n" + utils.CURRENT_TIMESTAMP + "\n" + Cfg.General.Name +
			"\n\nINFO" +
			"\nFlow update request    : " + strconv.Itoa(TDSStatistic.UpdatedFlows) +
			"\nFlow appended          : " + strconv.Itoa(TDSStatistic.AppendedFlows) +
			// TODO: для имплементации кода пикселя
			//"\nPixel request          : " + strconv.Itoa(TDSStatistic.PixelRequest) +
			"\nClick build request    : " + humanize.Comma(int64(TDSStatistic.ClickBuildRequest)) +
			"\nFlow Info request      : " + humanize.Comma(int64(TDSStatistic.FlowInfoRequest)) +
			"\nRedirect request       : " + humanize.Comma(int64(TDSStatistic.RedirectRequest)) +
			"\nIncorrect request      : " + humanize.Comma(int64(TDSStatistic.IncorrectRequest)) +
			"\nCookies request        : " + humanize.Comma(int64(TDSStatistic.CookieRequest)) +
			"\nUnique request (?)     : " + humanize.Comma(int64(uniqueRequests)) +
			"\n\nSYSTEM" +
			"\nOperating system       : " + Cfg.General.OS +
			"\nDebug level            : " + strconv.Itoa(Cfg.Debug.Level) +
			"\nTotal memory allocated : " + memoryUsageGeneral + " mb" +
			"\nPrivate memory         : " + memoryUsagePrivate + " mb" +
			"\nOpened files           : " + openedFiles +
			"\nCurrent rate           : " + humanize.Comma(int64(currentRPS)) + " rps" +
			"\nAverage rate           : " + humanize.Comma(int64(averageRPS)) + " rps" +
			"\nUptime                 : " + uptime +
			"\nProcessing time        : " + processingTime +
			"\nAverage response time  : " + avgReq +
			"\n\nREDIS" +
			"\nConnection             : " + strconv.FormatBool(IsRedisAlive) +
			"\nClicks sent/saved      : " + humanize.Comma(int64(TDSStatistic.ClicksSentToRedis)) +
			"\n"

		// добавим ++ к счетчику
		_ = Redisdb.Incr("StatisticCounter").Err()

		return text
	} else {
		TDSStatistic.Reset()

		//setting counter up
		_, err := Redisdb.Get("StatisticCounter").Result()
		if err != nil {
			_ = Redisdb.Set("StatisticCounter", 0, 0).Err()
		}

		return text
	}
}
