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
	m.Alloc = BToMb(rtm.Alloc)
	m.TotalAlloc = BToMb(rtm.TotalAlloc)
	m.Sys = BToMb(rtm.Sys)
	m.Mallocs = BToMb(rtm.Mallocs)
	m.Frees = BToMb(rtm.Frees)

	m.HeapSys = BToMb(rtm.HeapSys)

	m.HeapAlloc = BToMb(rtm.HeapAlloc)
	m.HeapIdle = BToMb(rtm.HeapIdle)
	m.HeapInuse = BToMb(rtm.HeapInuse)

	m.RealDetected = BToMb(rtm.Sys) + BToMb(rtm.HeapSys) + BToMb(rtm.HeapAlloc) + BToMb(rtm.HeapInuse) - BToMb(rtm.Alloc)
	m.RealDetected2 = BToMb(rtm.HeapSys) - BToMb(rtm.Alloc)

	m.HeapReleased = BToMb(rtm.HeapReleased)
	m.HeapObjects = BToMb(rtm.HeapObjects)

	m.StackInuse = BToMb(rtm.StackInuse)
	m.StackSys = BToMb(rtm.StackSys)
	m.MSpanInuse = BToMb(rtm.MSpanInuse)
	m.MSpanSys = BToMb(rtm.MSpanSys)
	m.MCacheInuse = BToMb(rtm.MCacheInuse)
	m.MCacheSys = BToMb(rtm.MCacheSys)
	m.BuckHashSys = BToMb(rtm.BuckHashSys)
	m.GCSys = BToMb(rtm.GCSys)
	m.OtherSys = BToMb(rtm.OtherSys)

	// Live objects = Mallocs - Frees
	m.LiveObjects = m.Mallocs - m.Frees

	// GC Stats
	m.PauseTotalNs = rtm.PauseTotalNs
	m.NumGC = rtm.NumGC

	return m
}
