package main

import (
	"strconv"
)

const flowModuleName = "flow.go"

//Prelands mini struct
type Prelands struct {
	ID  string
	URL string
}

//Lands mini struct
type Lands struct {
	ID  string
	URL string
}

//Counters mini struct
type Counters struct {
	ID   string
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

func (Flow FlowData) getInfo(FlowHash string) FlowData {
	//Фллоу хеш, если он есть в базе значит все пучком
	FlowID, _ := redisdb.Get(FlowHash + ":ID").Result()

	Flow.Hash = FlowHash
	convertedID, _ := strconv.Atoi(FlowID)
	Flow.ID = convertedID

	FlowWebMasterID, _ := redisdb.Get(FlowHash + ":WebMasterID").Result()
	convertedID, _ = strconv.Atoi(FlowWebMasterID)
	Flow.WebMasterID = convertedID

	FlowWebMasterCurrencyID, _ := redisdb.Get(FlowHash + ":WebMasterCurrencyID").Result()
	convertedID, _ = strconv.Atoi(FlowWebMasterCurrencyID)
	Flow.WebMasterCurrencyID = convertedID

	FlowOfferID, _ := redisdb.Get(FlowHash + ":OfferID").Result()
	convertedID, _ = strconv.Atoi(FlowOfferID)
	Flow.OfferID = convertedID

	// список лендингов
	FlowLandsList, _ := redisdb.Keys(FlowHash + ":land:*").Result()
	for _, item := range FlowLandsList {
		dataID, _ := redisdb.HGet(item, "id").Result()
		dataURL, _ := redisdb.HGet(item, "url").Result()
		Flow.Lands = append(Flow.Lands, Lands{dataID, dataURL})
	}

	FlowPrelandsList, _ := redisdb.Keys(FlowHash + ":preland:*").Result()
	for _, item := range FlowPrelandsList {
		dataID, _ := redisdb.HGet(item, "id").Result()
		dataURL, _ := redisdb.HGet(item, "url").Result()
		Flow.Prelands = append(Flow.Prelands, Prelands{dataID, dataURL})
	}

	FlowCountersList, _ := redisdb.Keys(FlowHash + ":counters:*").Result()
	for _, item := range FlowCountersList {
		dataID, _ := redisdb.HGet(item, "id").Result()
		dataName, _ := redisdb.HGet(item, "name").Result()
		Flow.Counters = append(Flow.Counters, Counters{dataID, dataName})
	}

	return Flow
}
