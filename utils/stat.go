package utils

import "time"

type TDSStats struct {
	UpdatedFlows     int
	AppendedFlows    int
	PixelRequest     int
	ClickInfoRequest int
	FlowInfoRequest  int
	RedirectRequest  int
	RedisStatRequest int
	IncorrectRequest int
	WorkTime         time.Duration
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
	tdstat.WorkTime = 0
}
