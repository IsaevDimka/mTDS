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
	"metatds/stuff/other/utils"
	"runtime"
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
	RealDetected2,
	RealDetected uint64
}

func MemMonitor() Monitor {
	var m Monitor
	var rtm runtime.MemStats
	// Read full mem stats
	runtime.ReadMemStats(&rtm)

	// Number of goroutines
	m.NumGoroutine = runtime.NumGoroutine()

	// Misc memory stats
	m.Alloc = utils.BToMb(rtm.Alloc)
	m.TotalAlloc = utils.BToMb(rtm.TotalAlloc)
	m.Sys = utils.BToMb(rtm.Sys)
	m.Mallocs = utils.BToMb(rtm.Mallocs)
	m.Frees = utils.BToMb(rtm.Frees)

	m.HeapSys = utils.BToMb(rtm.HeapSys)

	m.HeapAlloc = utils.BToMb(rtm.HeapAlloc)
	m.HeapIdle = utils.BToMb(rtm.HeapIdle)
	m.HeapInuse = utils.BToMb(rtm.HeapInuse)

	m.RealDetected = utils.BToMb(rtm.Sys) + utils.BToMb(rtm.HeapSys) + utils.BToMb(rtm.HeapAlloc) + utils.BToMb(rtm.HeapInuse) - utils.BToMb(rtm.Alloc)
	m.RealDetected2 = utils.BToMb(rtm.HeapSys) - utils.BToMb(rtm.Alloc)

	m.HeapReleased = utils.BToMb(rtm.HeapReleased)
	m.HeapObjects = utils.BToMb(rtm.HeapObjects)

	m.StackInuse = utils.BToMb(rtm.StackInuse)
	m.StackSys = utils.BToMb(rtm.StackSys)
	m.MSpanInuse = utils.BToMb(rtm.MSpanInuse)
	m.MSpanSys = utils.BToMb(rtm.MSpanSys)
	m.MCacheInuse = utils.BToMb(rtm.MCacheInuse)
	m.MCacheSys = utils.BToMb(rtm.MCacheSys)
	m.BuckHashSys = utils.BToMb(rtm.BuckHashSys)
	m.GCSys = utils.BToMb(rtm.GCSys)
	m.OtherSys = utils.BToMb(rtm.OtherSys)

	// Live objects = Mallocs - Frees
	m.LiveObjects = m.Mallocs - m.Frees

	// GC Stats
	m.PauseTotalNs = rtm.PauseTotalNs
	m.NumGC = rtm.NumGC

	return m
}
