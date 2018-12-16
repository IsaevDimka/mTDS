package utils

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/proxy"
)

type TelegramAdapter struct {
	Chats        []string
	User         string
	Password     string
	ProxyAddress string
	URL          string
	Token        string
	isActive     bool
}

func (tg *TelegramAdapter) Init(Ch []string, User, Password, Proxy, SendURL, SendToken string) bool {
	if len(Ch) > 0 && User != "" && Password != "" && Proxy != "" && SendURL != "" && SendToken != "" {
		tg.Chats = Ch
		tg.User = User
		tg.Password = Password
		tg.ProxyAddress = Proxy
		tg.URL = SendURL
		tg.Token = SendToken
		tg.isActive = true
		return true
	}
	return false
}

func (tg *TelegramAdapter) SendMessage(text string) bool {
	if tg.isActive {

		var authorization = new(proxy.Auth)
		authorization.User = tg.User
		authorization.Password = tg.Password

		// create a socks5 dialer
		dialer, err := proxy.SOCKS5("tcp", tg.ProxyAddress, authorization, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			//	os.Exit(1)
		}
		// setup a http client
		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}
		// set our socks5 as the dialer
		httpTransport.Dial = dialer.Dial
		// create a request

		for _, item := range tg.Chats {
			req, err := http.NewRequest("GET", tg.URL+tg.Token+"/sendMessage?chat_id="+item+"&text="+text, nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, "can't create request:", err)
				//	os.Exit(2)
			}
			// use the http client to fetch the page
			resp, err := httpClient.Do(req)
			if err != nil {
				fmt.Fprintln(os.Stderr, "can't GET page:", err)
				//	os.Exit(3)
			}
			defer resp.Body.Close()
		}
		//b, err := ioutil.ReadAll(resp.Body)
		//if err != nil {
		//	fmt.Fprintln(os.Stderr, "error reading body:", err)
		//	os.Exit(4)
		//}
		//fmt.Println(string(b))
		return true
	}
	return false
}

// const (
// 	PROXY_ADDR = "185.244.219.200:5443"
// 	URL        = "https://api.telegram.org/bot"
// 	TOKEN      = "624727745:AAH2YZUaw7OCkpch76RqH4ciggDmy9PvuoI"
// 	PROXY_USER = "mihail"
// 	PROXY_PASSWORD = "fit19280614"
// )
//
// var chats = []string{"232291795"}

// func main() {
// 	var tg TelegramAdapter
// 	tg.Init(chats,PROXY_USER,PROXY_PASSWORD,PROXY_ADDR,URL,TOKEN)
//
// 	if !tg.SendMessage("shit") {
// 		fmt.Println("Error sending message")
// 	}
// }
//
