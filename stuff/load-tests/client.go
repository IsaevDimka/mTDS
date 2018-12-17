package main

import ("fmt";"time";"net/http";"io/ioutil")

func main(){

	start:=time.Now()

	for i:=0; i<3000; i++ {
		resp, err := http.Get("http://localhost/")
		if err != nil {
			// handle error
		} 
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)

		fmt.Println("Query #",i)
	}

	fmt.Println("[ ELAPSED TIME ]",time.Since(start))

}


