package utils

import "time"

var ResponseAverage []time.Duration

type TDSStats struct {
	UpdatedFlows     int
	AppendedFlows    int
	PixelRequest     int
	ClickInfoRequest int
	FlowInfoRequest  int
	RedirectRequest  int
	RedisStatRequest int
	IncorrectRequest int
	CookieRequest    int
	ProcessingTime   time.Duration
	MemoryAllocated  uint64
}

func (tdstat TDSStats) Reset() {
	tdstat.UpdatedFlows = 0
	tdstat.AppendedFlows = 0
	tdstat.PixelRequest = 0
	tdstat.ClickInfoRequest = 0
	tdstat.RedisStatRequest = 0
	tdstat.FlowInfoRequest = 0
	tdstat.RedirectRequest = 0
	tdstat.IncorrectRequest = 0
	tdstat.CookieRequest = 0
	tdstat.ProcessingTime = 0
	tdstat.MemoryAllocated = 0

	ResponseAverage = []time.Duration{}
}
