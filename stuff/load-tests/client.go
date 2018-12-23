package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {

	start := time.Now()

	for i := 0; i < 3000; i++ {
		resp, err := http.Get("http://localhost/001kezbiav/1/2/3/4/5")
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)

		fmt.Println("Query #", i)
	}

	fmt.Println("[ ELAPSED TIME ]", time.Since(start))

}
