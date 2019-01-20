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
	Name string
	ID   int
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

type FlowImportData struct {
	ID                  int
	OfferID             int
	WebMasterID         int
	WebMasterCurrencyID int
	Hash                string
	Prelands            []Prelands
	Lands               []Lands
	Counters            []Counters
}

func (Flow FlowData) GetInfo(FlowHash string) FlowData {
	//Фллоу хеш, если он есть в базе значит все пучком
	FlowResult, _ := config.Redisdb.HGetAll(FlowHash).Result()

	Flow.Hash = FlowResult["Hash"]
	convertedID, _ := strconv.Atoi(FlowResult["ID"])
	Flow.ID = convertedID
	convertedID, _ = strconv.Atoi(FlowResult["WebMasterID"])
	Flow.WebMasterID = convertedID
	convertedID, _ = strconv.Atoi(FlowResult["WebMasterCurrencyID"])
	Flow.WebMasterCurrencyID = convertedID
	convertedID, _ = strconv.Atoi(FlowResult["OfferID"])
	Flow.OfferID = convertedID

	// получаем массивы
	FlowLandsList, _ := config.Redisdb.HGetAll(FlowHash + ":lands").Result()
	for i, item := range FlowLandsList {
		convertedID, _ = strconv.Atoi(i)
		Flow.Lands = append(Flow.Lands, Lands{convertedID, item})
	}

	FlowPrelandsList, _ := config.Redisdb.HGetAll(FlowHash + ":prelands").Result()
	for i, item := range FlowPrelandsList {
		convertedID, _ = strconv.Atoi(i)
		Flow.Prelands = append(Flow.Prelands, Prelands{convertedID, item})
	}

	FlowCountersList, _ := config.Redisdb.HGetAll(FlowHash + ":counters").Result()
	for i, item := range FlowCountersList {
		convertedID, _ = strconv.Atoi(i)
		Flow.Counters = append(Flow.Counters, Counters{item, convertedID})
	}

	return Flow
}
