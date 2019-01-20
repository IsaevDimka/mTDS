/****************************************************************************************************
*
* Perks to use in other modules like system stat, special for Meta CPA, Ltd.
* by Michael S. Merzlyakov AFKA predator_pc@12122018
* version v2.0.3
*
* created at 04122018
* last edit: 16122018
*
*****************************************************************************************************/

package config

import (
	"github.com/davecgh/go-spew/spew"
	"time"
)

// Template for Channel by predator_pc
//
// func ChannelWithSleepTemplate() <-chan string {
// 	c := make(chan string)
// 	go func() {
// 		for {
// 			time.Sleep(time.Minute * 10)
// 		}
// 	}()
// 	return c
// }

//
// text representation with dump of variable \t \n delemiters
//
func GetSystemConfiguration() string {
	text := spew.Sdump(Cfg)
	return text
}

//
// average responder for stat by time of execution
//
func DurationAverage(dur []time.Duration) time.Duration {
	var allTime float64
	for _, item := range dur {
		allTime += float64(item)
	}
	result := allTime / float64(1+len(dur))
	return time.Duration(result)
}

//
// average responder for stat by time of execution
//
func RPSAverage(RPSStat []int) int {
	var allStat float64
	for _, item := range RPSStat {
		allStat += float64(item)
	}
	result := allStat / float64(1+len(RPSStat))
	return int(result)
}
