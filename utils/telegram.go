/****************************************************************************************************
*
* Telegram Adapter module, special for Meta CPA, Ltd.
* by Michael S. Merzlyakov AFKA predator_pc@16122018
* version v2.0.3
*
* created at 14122018
* last edit: 16122018
*
*****************************************************************************************************/

package utils

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type TelegramAdapter struct {
	Chats        []string
	User         string
	Password     string
	ProxyAddress string
	URL          string
	Token        string
	UseProxy     bool
	isActive     bool
}

func (tg *TelegramAdapter) Init(Ch []string, User, Password, Proxy, SendURL, SendToken string, UseProxy bool) bool {
	if len(Ch) > 0 && User != "" && Password != "" && Proxy != "" && SendURL != "" &&
		SendToken != "" {
		tg.Chats = Ch
		tg.User = User
		tg.Password = Password
		tg.ProxyAddress = Proxy
		tg.URL = SendURL
		tg.Token = SendToken
		tg.isActive = true
		tg.UseProxy = UseProxy
		return true
	}
	return false
}

func (tg *TelegramAdapter) SendMessage(text string) bool {
	if tg.isActive {

		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}

		if  tg.UseProxy {
			// setup a http client
			var authorization= new(proxy.Auth)
			authorization.User = tg.User
			authorization.Password = tg.Password

			// create a socks5 dialer
			dialer, err := proxy.SOCKS5("tcp", tg.ProxyAddress, authorization, proxy.Direct)
			if err != nil {
				fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			}

			// set our socks5 as the dialer
			httpTransport.Dial = dialer.Dial
		}

		for _, item := range tg.Chats {
			if !tg.UseProxy {
				req, err := http.Get(tg.URL+tg.Token+"/sendMessage?parse_mode=markdown&chat_id="+item+"&text="+url.QueryEscape(text))

				if req != nil {
					defer req.Body.Close()
					ioutil.ReadAll(req.Body)
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, "can't create request:", err)
				}
			} else {
				req, err := http.NewRequest("GET", tg.URL+tg.Token+"/sendMessage?parse_mode=markdown&chat_id="+item+"&text="+url.QueryEscape(text), nil)

				// use the http client to fetch the page
				resp, err := httpClient.Do(req)
				if err != nil {
					// TODO this needs to be recovered from panic otherwise fails
					fmt.Fprintln(os.Stderr, "can't GET page:", err)
				}
				defer resp.Body.Close()
			}
		}
		return true
	}
	return false
}
