package main

import (
	"strconv"
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
	Location            string
	Sub1                string
	Sub2                string
	Sub3                string
	Sub4                string
	Sub5                string
}

func (ClickData) generateCID() string {
	rnd := randomString(cfg.Click.Length) //generateRandomString(cfg.Click.Length)
	// TODO: Здесь надо покрутить проверку, на то, что сейчас в редисе на предмет дублирующих
	// сидов, иначе жопа может быть
	if cfg.Debug.Level > 1 {
		printInfo("Click ID", rnd, clickModuleName)
	}
	return rnd
}

/*
	saving click to redis
*/

func (Click ClickData) save() bool {
	errFlowID := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "FlowID", Click.FlowID).Err()
	errLandingID := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "LandingID", Click.LandingID).Err()
	errPrelandingID := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "PrelandingID", Click.PrelandingID).Err()
	errWebMasterID := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "WebMasterID", Click.WebMasterID).Err()
	errWebMasterCurrencyID := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "WebMasterCurrencyID", Click.WebMasterCurrencyID).Err()
	errOfferID := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "OfferID", Click.OfferID).Err()

	errFlowHash := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "FlowHash", Click.FlowHash).Err()
	errHash := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Hash", Click.Hash).Err()
	errReferer := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Referer", Click.Referer).Err()
	errTime := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Time", Click.Time).Err()
	errIP := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "IP", Click.IP).Err()
	errUserAgent := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "UserAgent", Click.UserAgent).Err()
	errURL := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "URL", Click.URL).Err()
	errLocation := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Location", Click.Location).Err()

	errSub1 := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Sub1", Click.Sub1).Err()
	errSub2 := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Sub2", Click.Sub2).Err()
	errSub3 := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Sub3", Click.Sub3).Err()
	errSub4 := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Sub4", Click.Sub4).Err()
	errSub5 := redisdb.HSet(Click.FlowHash+":click:"+Click.Hash, "Sub5", Click.Sub5).Err()

	// TODO: надо придумать повеселее
	if errFlowID != nil || errLandingID != nil || errPrelandingID != nil || errWebMasterID != nil ||
		errWebMasterCurrencyID != nil || errOfferID != nil || errFlowHash != nil || errHash != nil ||
		errReferer != nil || errTime != nil || errIP != nil || errUserAgent != nil || errURL != nil ||
		errLocation != nil || errSub1 != nil || errSub2 != nil || errSub3 != nil || errSub4 != nil || errSub5 != nil {
		return false
	} else {
		return true
	}
}

func (Click ClickData) getInfo(ClickHash string) ClickData {
	//var Click ClickData
	var FlowHash string

	Click.Hash = ClickHash
	FlowHashKeys, _ := redisdb.Keys("*:click:" + Click.Hash).Result()

	if len(FlowHashKeys) > 0 {
		FlowHash, _ = redisdb.HGet(FlowHashKeys[0], "FlowHash").Result()
	}

	// int should be converted
	Click.FlowHash = FlowHash
	ClickFlowID, _ := redisdb.HGet(Click.FlowHash+":click:"+ClickHash, "FlowID").Result()

	convertedID, _ := strconv.Atoi(ClickFlowID)
	Click.FlowID = convertedID

	ClickLandingID, _ := redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "LandingID").Result()
	convertedID, _ = strconv.Atoi(ClickLandingID)
	Click.LandingID = convertedID

	ClickPrelandingID, _ := redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "PrelandingID").Result()
	convertedID, _ = strconv.Atoi(ClickPrelandingID)
	Click.PrelandingID = convertedID

	ClickWebMasterID, _ := redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "WebMasterID").Result()
	convertedID, _ = strconv.Atoi(ClickWebMasterID)
	Click.WebMasterID = convertedID

	ClickWebMasterCurrencyID, _ := redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "WebMasterCurrencyID").Result()
	convertedID, _ = strconv.Atoi(ClickWebMasterCurrencyID)
	Click.WebMasterCurrencyID = convertedID

	ClickOfferID, _ := redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "OfferID").Result()
	convertedID, _ = strconv.Atoi(ClickOfferID)
	Click.OfferID = convertedID

	Click.Hash, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Hash").Result()
	Click.Referer, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Referer").Result()
	Click.Time, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Time").Result()
	Click.IP, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "IP").Result()
	Click.UserAgent, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "UserAgent").Result()
	Click.URL, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "URL").Result()
	Click.Location, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Location").Result()

	// subs
	Click.Sub1, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Sub1").Result()
	Click.Sub2, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Sub2").Result()
	Click.Sub3, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Sub3").Result()
	Click.Sub4, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Sub4").Result()
	Click.Sub5, _ = redisdb.HGet(Click.FlowHash+":click:"+Click.Hash, "Sub5").Result()

	return Click
}
