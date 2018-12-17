package main

import ("fmt";u "pack/utils";"time";"net/http";"net/url";"log";"io/ioutil")

func main(){
	start:= time.Now()

	var v   u.Utils
	v.ID    = 10
	v.Name  = "Jack"	
	v.Some1 = "123"
	v.Some2 = "234"

	fmt.Println("[ Utils says ]",v.SayID(),v.SayName())

	for i:=0; i<10000000; i++ { 
	// 
	}

	fmt.Println("[ Time elsapsed ]",time.Since(start))

//	resp, err := 
	http.PostForm("https://api.telegram.org/bot"+"624727745:AAH2YZUaw7OCkpch76RqH4ciggDmy9PvuoI"+ "/sendMessage?chat_id=232291795"+"&text=Yes, motherfucker!",url.Values{})
//	url.Values{"key": {"Value"}, "id": {"123"}})


//	client := &http.Client{}
//	client.SetProxy("167.99.74.231:3128") 
//	os.Setenv("HTTP_PROXY", "https://167.99.74.231:3128")
//	client.Get("https://api.telegram.org/bot"+"624727745:AAH2YZUaw7OCkpch76RqH4ciggDmy9PvuoI"+ "/sendMessage?chat_id=232291795"+"&text=Yes, motherfucker!") // do request through proxy

//	request := gorequest.New().Proxy("http://138.219.229.239:80")
//	request.Get("https://api.telegram.org/bot"+"624727745:AAH2YZUaw7OCkpch76RqH4ciggDmy9PvuoI"+ "/sendMessage?chat_id=232291795"+"&text=Yesmotherfucker!").End()

  //creating the proxyURL
    proxyStr := "http://121.52.157.23:8080"
    proxyURL, err := url.Parse(proxyStr)
    if err != nil {
        log.Println(err)
    }

    //creating the URL to be loaded through the proxy
    urlStr := "https://api.telegram.org/bot"+"624727745:AAH2YZUaw7OCkpch76RqH4ciggDmy9PvuoI"+ "/sendMessage?chat_id=232291795"+"&text=Yesmotherfucker!"
    url, err := url.Parse(urlStr)
    if err != nil {
        log.Println(err)
    }

    //adding the proxy settings to the Transport object
    transport := &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    }

    //adding the Transport object to the http Client
    client := &http.Client{
        Transport: transport,
    }

    //generating the HTTP GET request
    request, err := http.NewRequest("GET", url.String(), nil)
    if err != nil {
        log.Println(err)
    }

    //calling the URL
    response, err := client.Do(request)
    if err != nil {
        log.Println(err)
    }

    //getting the response
    data, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Println(err)
    }
    //printing the response
    log.Println(string(data))

}


