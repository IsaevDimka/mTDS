/****************************************************************************************************
*
* Click model/entity and methods, special for Meta CPA, Ltd.
* by Michael S. Merzlyakov AFKA predator_pc@12122018
* version v2.0.3
*
* created at 04122018
* last edit: 16122018
*
*****************************************************************************************************/

package models

import (
	"metatds/config"
	"metatds/utils"
	"strconv"
	"github.com/fatih/structs"
)

const clickModuleName = "click.go"

//ClickData for click info
type ClickData struct {
	FlowID              int
	LandingID           int
	PrelandingID        int
	OfferID             int
	WebMasterID         int
	WebMasterCurrencyID int
	FlowHash            string
	Hash                string
	IP                  string
	URL                 string
	Time                string
	Referer             string
	UserAgent           string
	LocationLP          string
	IsVisitedLP         int
	LocationPL          string
	IsVisitedPL         int
	Sub1                string
	Sub2                string
	Sub3                string
	Sub4                string
	Sub5                string
}

func (ClickData) GenerateCID() string {
//restart:
	rnd := utils.RandomString(config.Cfg.Click.Length)
	return rnd
	//
	// if config.Cfg.Debug.Level > 1 {
	// 	utils.PrintInfo("Click ID", rnd, clickModuleName)
	// }
	//
	// if clickHashExists(rnd) {
	// 	utils.PrintError("Click ID already exists", rnd, clickModuleName)
	// 	goto restart
	// } else {
	// 	return rnd
	// }
}

func (Click ClickData) Save() bool {
	c:=structs.Map(Click)
	if err := config.Redisdb.HMSet(Click.FlowHash+":click:"+Click.Hash, c).Err(); err!=nil {
		return false
	}
	return true
}

//
// /*
// 	saving click to redis
// */
// func (Click ClickData) Save() bool {
// 	errFlowID := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "FlowID", Click.FlowID).Err()
// 	errLandingID := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "LandingID", Click.LandingID).Err()
// 	errPrelandingID := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "PrelandingID", Click.PrelandingID).Err()
// 	errWebMasterID := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "WebMasterID", Click.WebMasterID).Err()
// 	errWebMasterCurrencyID := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "WebMasterCurrencyID", Click.WebMasterCurrencyID).Err()
// 	errOfferID := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "OfferID", Click.OfferID).Err()
//
// 	errFlowHash := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "FlowHash", Click.FlowHash).Err()
// 	errHash := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Hash", Click.Hash).Err()
// 	errReferer := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Referer", Click.Referer).Err()
// 	errTime := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Time", Click.Time).Err()
// 	errIP := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "IP", Click.IP).Err()
// 	errUserAgent := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "UserAgent", Click.UserAgent).Err()
// 	errURL := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "URL", Click.URL).Err()
// 	errLocationLP := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "LocationLP", Click.LocationLP).Err()
// 	errIsVisitedLP := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "IsVisitedLP", Click.IsVisitedLP).Err()
//
// 	errLocationPL := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "LocationPL", Click.LocationPL).Err()
// 	errIsVisitedPL := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "IsVisitedPL", Click.IsVisitedPL).Err()
//
// 	errSub1 := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Sub1", Click.Sub1).Err()
// 	errSub2 := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Sub2", Click.Sub2).Err()
// 	errSub3 := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Sub3", Click.Sub3).Err()
// 	errSub4 := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Sub4", Click.Sub4).Err()
// 	errSub5 := config.Redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Sub5", Click.Sub5).Err()
//
// 	// TODO: надо придумать повеселее
// 	if errFlowID != nil || errLandingID != nil || errPrelandingID != nil || errWebMasterID != nil ||
// 		errWebMasterCurrencyID != nil || errOfferID != nil || errFlowHash != nil || errHash != nil ||
// 		errReferer != nil || errTime != nil || errIP != nil || errUserAgent != nil || errURL != nil ||
// 		errLocationLP != nil || errLocationPL != nil || errIsVisitedLP != nil || errIsVisitedPL != nil ||
// 		errSub1 != nil || errSub2 != nil || errSub3 != nil || errSub4 != nil ||
// 		errSub5 != nil {
// 		return false
// 	} else {
// 		return true
// 	}
// }

func (Click ClickData) GetInfo(ClickHash string) ClickData {

	//
	// TODO Maybe it would be better to use HGetAll instead of Keys/HGet
	//

	var FlowHash string

	Click.Hash = ClickHash
	FlowHashKeys, _ := config.Redisdb.Keys("*:click:" + Click.Hash).Result()

	if len(FlowHashKeys) > 0 {
		FlowHash, _ = config.Redisdb.HGet(FlowHashKeys[0], "FlowHash").Result()
	}

	// int should be converted
	Click.FlowHash = FlowHash
	ClickFlowID, _ := config.Redisdb.HGet(Click.FlowHash+":click:"+ClickHash, "FlowID").Result()

	convertedID, _ := strconv.Atoi(ClickFlowID)
	Click.FlowID = convertedID

	ClickLandingID, _ := config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "LandingID").Result()
	convertedID, _ = strconv.Atoi(ClickLandingID)
	Click.LandingID = convertedID

	ClickPrelandingID, _ := config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "PrelandingID").Result()
	convertedID, _ = strconv.Atoi(ClickPrelandingID)
	Click.PrelandingID = convertedID

	ClickWebMasterID, _ := config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "WebMasterID").Result()
	convertedID, _ = strconv.Atoi(ClickWebMasterID)
	Click.WebMasterID = convertedID

	ClickWebMasterCurrencyID, _ := config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "WebMasterCurrencyID").Result()
	convertedID, _ = strconv.Atoi(ClickWebMasterCurrencyID)
	Click.WebMasterCurrencyID = convertedID

	ClickOfferID, _ := config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "OfferID").Result()
	convertedID, _ = strconv.Atoi(ClickOfferID)
	Click.OfferID = convertedID

	Click.Hash, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Hash").Result()
	Click.Referer, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Referer").Result()
	Click.Time, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Time").Result()
	Click.IP, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "IP").Result()
	Click.UserAgent, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "UserAgent").Result()
	Click.URL, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "URL").Result()
	Click.LocationLP, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "LocationLP").Result()
	ClickIsVisitedLP, _ := config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "IsVisitedLP").Result()
	convertedID, _ = strconv.Atoi(ClickIsVisitedLP)
	Click.IsVisitedLP = convertedID

	Click.LocationPL, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "LocationPL").Result()
	ClickIsVisitedPL, _ := config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "IsVisitedPL").Result()
	convertedID, _ = strconv.Atoi(ClickIsVisitedPL)
	Click.IsVisitedPL = convertedID

	// subs
	Click.Sub1, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Sub1").Result()
	Click.Sub2, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Sub2").Result()
	Click.Sub3, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Sub3").Result()
	Click.Sub4, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Sub4").Result()
	Click.Sub5, _ = config.Redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Sub5").Result()

	return Click
}

func clickHashExists(ClickHash string) bool {
	keys, err := config.Redisdb.Keys("*:click:" + ClickHash).Result()
	if err != nil || len(keys) > 0 {
		return true
	}
	return false
}
