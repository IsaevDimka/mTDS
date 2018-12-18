/****************************************************************************************************
*
* Flow model/entity and methods, special for Meta CPA, Ltd.
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
	"strconv"
)

const flowModuleName = "flow.go"

//Prelands mini struct
type Prelands struct {
	ID  int
	URL string
}

//Lands mini struct
type Lands struct {
	ID  int
	URL string
}

//Counters mini struct
type Counters struct {
	ID   int
	Name string
}

//FlowData for flow info
type FlowData struct {
	ID                  int
	OfferID             int
	WebMasterID         int
	WebMasterCurrencyID int
	Hash                string
	RandomPreland       string
	RandomLand          string
	Prelands            []Prelands
	Lands               []Lands
	Counters            []Counters
}

func (Flow FlowData) GetInfo(FlowHash string) FlowData {
	//Фллоу хеш, если он есть в базе значит все пучком
	FlowID, _ := config.Redisdb.Get(FlowHash + ":ID").Result()

	Flow.Hash = FlowHash
	convertedID, _ := strconv.Atoi(FlowID)
	Flow.ID = convertedID

	FlowWebMasterID, _ := config.Redisdb.Get(FlowHash + ":WebMasterID").Result()
	convertedID, _ = strconv.Atoi(FlowWebMasterID)
	Flow.WebMasterID = convertedID

	FlowWebMasterCurrencyID, _ := config.Redisdb.Get(FlowHash + ":WebMasterCurrencyID").Result()
	convertedID, _ = strconv.Atoi(FlowWebMasterCurrencyID)
	Flow.WebMasterCurrencyID = convertedID

	FlowOfferID, _ := config.Redisdb.Get(FlowHash + ":OfferID").Result()
	convertedID, _ = strconv.Atoi(FlowOfferID)
	Flow.OfferID = convertedID

	// список лендингов
	FlowLandsList, _ := config.Redisdb.Keys(FlowHash + ":land:*").Result()
	for _, item := range FlowLandsList {
		dataID, _ := config.Redisdb.HGet(item, "id").Result()
		dataURL, _ := config.Redisdb.HGet(item, "url").Result()
		convertedID, _ = strconv.Atoi(dataID)
		Flow.Lands = append(Flow.Lands, Lands{convertedID, dataURL})
	}

	FlowPrelandsList, _ := config.Redisdb.Keys(FlowHash + ":preland:*").Result()
	for _, item := range FlowPrelandsList {
		dataID, _ := config.Redisdb.HGet(item, "id").Result()
		dataURL, _ := config.Redisdb.HGet(item, "url").Result()
		convertedID, _ = strconv.Atoi(dataID)
		Flow.Prelands = append(Flow.Prelands, Prelands{convertedID, dataURL})
	}

	FlowCountersList, _ := config.Redisdb.Keys(FlowHash + ":counters:*").Result()
	for _, item := range FlowCountersList {
		dataID, _ := config.Redisdb.HGet(item, "id").Result()
		dataName, _ := config.Redisdb.HGet(item, "name").Result()
		convertedID, _ = strconv.Atoi(dataID)
		Flow.Counters = append(Flow.Counters, Counters{convertedID, dataName})
	}

	return Flow
}
