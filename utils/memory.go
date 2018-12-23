/****************************************************************************************************
*
* Memory and system monitoring module, special for Meta CPA, Ltd.
* by Michael S. Merzlyakov AFKA predator_pc@06122018
* version v2.0.3
* special thanks to https://scene-si.org/2018/08/06/basic-monitoring-of-go-apps-with-the-runtime-package/
* created at 06122018
* last edit: 16122018
*
*****************************************************************************************************/

package utils

import (
	"fmt"
	"runtime"
	"time"
)

type Monitor struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	HeapSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapReleased,
	HeapObjects,
	StackInuse,
	StackSys,
	MSpanInuse,
	MSpanSys,
	MCacheInuse,
	MCacheSys,
	BuckHashSys,
	GCSys,
	OtherSys,


	PauseTotalNs uint64

	NumGC        uint32
	NumGoroutine int
}

func MemMonitor(duration int) {
	var m Monitor
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)

		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		m.NumGoroutine = runtime.NumGoroutine()

		// Misc memory stats
		m.Alloc = BToKb(rtm.Alloc)
		m.TotalAlloc = BToKb(rtm.TotalAlloc)
		m.Sys = BToKb(rtm.Sys)
		m.Mallocs = BToKb(rtm.Mallocs)
		m.Frees = BToKb(rtm.Frees)

		m.HeapSys = BToKb(rtm.HeapSys)
		m.HeapAlloc = BToKb(rtm.HeapAlloc)
		m.HeapIdle = BToKb(rtm.HeapIdle)
		m.HeapInuse = BToKb(rtm.HeapInuse)
		m.HeapReleased = BToKb(rtm.HeapReleased)
		m.HeapObjects = BToKb(rtm.HeapObjects)


		m.StackInuse = BToKb(rtm.StackInuse)
		m.StackSys = BToKb(rtm.StackSys)
		m.MSpanInuse = BToKb(rtm.MSpanInuse)
		m.MSpanSys = BToKb(rtm.MSpanSys)
		m.MCacheInuse = BToKb(rtm.MCacheInuse)
		m.MCacheSys = BToKb(rtm.MCacheSys)
		m.BuckHashSys = BToKb(rtm.BuckHashSys)
		m.GCSys = BToKb(rtm.GCSys)
		m.OtherSys = BToKb(rtm.OtherSys)

		// Live objects = Mallocs - Frees
		m.LiveObjects = m.Mallocs - m.Frees

		// GC Stats
		m.PauseTotalNs = rtm.PauseTotalNs
		m.NumGC = rtm.NumGC

		// Just encode to json and print
		//b, _ := json.Marshal(m)
		fmt.Println(JSONPretty(m))
		//fmt.Printf("Total system memory: %d\n", memory.TotalMemory())
	}
}
