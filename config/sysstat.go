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

	StatisticCounter, err := Redisdb.Get("StatisticCounter").Result()
	if err != nil {
		_ = Redisdb.Set("StatisticCounter", 0, 0).Err()
		StatisticCounter = "0"
	}

	currentRPSstart := TDSStatistic.RedirectRequest
	time.Sleep(1 * time.Second)
	currentRPS := TDSStatistic.RedirectRequest - currentRPSstart

	if len(RPSStat) < minimumStatCountRPS {
		RPSStat = append(RPSStat, currentRPS)
	} else {
		RPSStat = RPSStatDefault
	}

	averageRPS := RPSAverage(RPSStat)

	if TDSStatistic != (utils.TDSStats{}) {
		duration = 60 * time.Minute

		if time.Since(UpTime) < duration {
			uptime = durafmt.Parse(time.Since(UpTime)).String(durafmt.DF_LONG)
			processingTime = durafmt.Parse(TDSStatistic.ProcessingTime).String(durafmt.DF_LONG)
		} else {
			uptime = durafmt.Parse(time.Since(UpTime)).String(durafmt.DF_MIDDLE)
			processingTime = durafmt.Parse(TDSStatistic.ProcessingTime).String(durafmt.DF_MIDDLE)
		}

		runtime.ReadMemStats(&memory)

		RealDetectedGeneral := memory.Sys + memory.HeapSys + memory.HeapAlloc + memory.HeapInuse - memory.Alloc
		RealDetectedPrivate := memory.HeapSys - memory.Alloc

		memoryUsageGeneral = strconv.FormatUint(utils.BToMb(RealDetectedGeneral), 10)
		memoryUsagePrivate = strconv.FormatUint(utils.BToMb(RealDetectedPrivate), 10)

		if Cfg.General.OS == "linux" || Cfg.General.OS == "unix" {

			pid := strconv.Itoa(os.Getpid())
			fds, e := ioutil.ReadDir("/proc/" + pid + "/fd")

			if e != nil && Cfg.Debug.Level > 0 {
				utils.PrintError("Error", "reading process directory failed", sysstatModuleName)
			} else {
				//utils.PrintInfo("PID", pid, initModuleName)
			}

			if len(fds) > 0 {
				openedFiles = strconv.Itoa(len(fds))
			}
		}

		dur := DurationAverage(utils.ResponseAverage)

		if dur < time.Duration(1*time.Millisecond) { //|| dur < time.Duration(1 * time.Microsecond) || dur < time.Duration(1 * time.Nanosecond) {
			avgReq = "< 1 msec"
		} else {
			avgReq = durafmt.Parse(dur).String(durafmt.DF_LONG)
		}

		uniqueRequests := (TDSStatistic.ClickInfoRequest + TDSStatistic.FlowInfoRequest + TDSStatistic.RedirectRequest) - TDSStatistic.CookieRequest // - TDSStatistic.IncorrectRequest

		//setting stat to redis
		_ = Redisdb.HSet("SystemStatistic", StatisticCounter, "["+StatisticCounter+","+
			fmt.Sprintf("%.0f", (math.Round(dur.Seconds()*1000)))+","+
			strconv.Itoa(averageRPS)+","+strconv.Itoa(currentRPS)+"]").Err()

		text = "\n" + utils.CURRENT_TIMESTAMP + "\n" + Cfg.General.Name +
			"\n\nINFO" +
			"\nFlow update request    : " + strconv.Itoa(TDSStatistic.UpdatedFlows) +
			"\nFlow appended          : " + strconv.Itoa(TDSStatistic.AppendedFlows) +
			//"\nPixel request          : " + strconv.Itoa(TDSStatistic.PixelRequest) +
			"\nClick Info request     : " + humanize.Comma(int64(TDSStatistic.ClickInfoRequest)) + //strconv.Itoa(TDSStatistic.ClickInfoRequest) +
			"\nFlow Info request      : " + humanize.Comma(int64(TDSStatistic.FlowInfoRequest)) + //strconv.Itoa(TDSStatistic.FlowInfoRequest) +
			"\nRedirect request       : " + humanize.Comma(int64(TDSStatistic.RedirectRequest)) + //strconv.Itoa(TDSStatistic.RedirectRequest) +
			//			"\nRedis Stat request     : " + strconv.Itoa(TDSStatistic.RedisStatRequest) +
			"\nIncorrect request      : " + humanize.Comma(int64(TDSStatistic.IncorrectRequest)) + //strconv.Itoa(TDSStatistic.IncorrectRequest) +
			"\nCookies request        : " + humanize.Comma(int64(TDSStatistic.CookieRequest)) + //strconv.Itoa(TDSStatistic.CookieRequest) +
			"\nUnique request (?)     : " + humanize.Comma(int64(uniqueRequests)) + //strconv.Itoa(uniqueRequests) +
			"\n\nSYSTEM" +
			"\nOperating system       : " + Cfg.General.OS +
			"\nDebug level            : " + strconv.Itoa(Cfg.Debug.Level) +
			"\nTotal memory allocated : " + memoryUsageGeneral + " mb" +
			"\nPrivate memory         : " + memoryUsagePrivate + " mb" +
			"\nOpened files           : " + openedFiles +
			"\nCurrent rate           : " + humanize.Comma(int64(currentRPS)) + " rps" +
			"\nAverage rate           : " + humanize.Comma(int64(averageRPS)) + " rps" + //strconv.Itoa(averageRPS) + " rps" +
			"\nUptime                 : " + uptime +
			"\nProcessing time        : " + processingTime +
			"\nAverage response time  : " + avgReq +
			"\n\nREDIS" +
			"\nConnection             : " + strconv.FormatBool(IsRedisAlive) +
			"\nClicks sent/saved      : " + humanize.Comma(int64(TDSStatistic.ClicksSentToRedis)) + //strconv.Itoa(TDSStatistic.ClicksSentToRedis) +
			"\n"

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
