package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"io/ioutil"
	"strconv"
	"time"
)

func JSONPretty(Data interface{}) string {
	var out bytes.Buffer //буфер конвертации джейсона в красивый джейсон
	jsonData, _ := json.Marshal(Data)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)
	json.Indent(&out, jsonData, "", "    ")
	return out.String()
}

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
	//RandomPreland       string
	//RandomLand          string
//	Prelands            []Prelands
	Lands               []Lands
	//Counters            []Counters
}

func main(){

	start:=time.Now()

	//get connection to Redis
	redisdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	b, err := ioutil.ReadFile("1.json")

	if err != nil {
		fmt.Print(err)
	}

	var data []FlowData

	if err := json.Unmarshal(b, &data); err != nil {
		panic(err)
	}

	//fmt.Println(JSONPretty(data))

	for _, item:= range data {
		_ = redisdb.Set(item.Hash+":ID",item.ID,0).Err()
		_ = redisdb.Set(item.Hash+":Hash",item.Hash,0).Err()
		_ = redisdb.Set(item.Hash+":OfferID",item.OfferID,0).Err()
		_ = redisdb.Set(item.Hash+":WebMasterID",item.WebMasterID,0).Err()
		_ = redisdb.Set(item.Hash+":WebMasterCurrencyID",item.WebMasterCurrencyID,0).Err()

		if len(item.Lands)>0 {
			for i, lands := range item.Lands {
				_ = redisdb.HSet(item.Hash + ":land:" + strconv.Itoa(i),"id",lands.ID)
				_ = redisdb.HSet(item.Hash + ":land:" + strconv.Itoa(i),"url",lands.URL)
			}
		}

	}
	//fmt.Println(item)


	fmt.Println("Time elapsed: ",time.Since(start))
}
