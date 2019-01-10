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
	"bytes"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net/http"
	"net/url"
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

		body := bytes.NewBuffer(nil)

		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}

		if tg.UseProxy {
			// setup a http client
			var authorization = new(proxy.Auth)
			authorization.User = tg.User
			authorization.Password = tg.Password

			// create a socks5 dialer
			dialer, err := proxy.SOCKS5("tcp", tg.ProxyAddress, authorization, proxy.Direct)
			if err != nil {
				fmt.Println("[ ERROR ] can't connect to the proxy: ", err)
			}

			// set our socks5 as the dialer
			httpTransport.Dial = dialer.Dial
		}

		markDownedText := "```" + text + "```"

		for _, item := range tg.Chats {
			if !tg.UseProxy {
				req, err := http.Get(tg.URL + tg.Token + "/sendMessage?parse_mode=markdown&chat_id=" + item + "&text=" + url.QueryEscape(markDownedText))

				if err != nil {
					recover()
					fmt.Println("[ ERROR] can't create request 0: ", err)
				}

				if req != nil {
					// setting header in case of API request is not NIL
					req.Header.Set("Connection", "close")
					// reading the body
					_, _ = io.Copy(body, req.Body)
					// closing anyway now
					// defer is not needed cause we get an exception before
					_ = req.Body.Close()
				}

			} else {
				req, err := http.NewRequest("GET", tg.URL+tg.Token+"/sendMessage?parse_mode=markdown&chat_id="+item+"&text="+url.QueryEscape(markDownedText), nil)

				if err != nil {
					recover()
					fmt.Println("[ ERROR] can't create request 1: ", err)
				} else {
					// we shouldn't write to request if doesn't created
					req.Header.Set("Connection", "close")
				}

				// use the http client to fetch the page
				resp, err := httpClient.Do(req)

				if err != nil {
					recover()
					fmt.Println("[ ERROR] can't create request 2: ", err)
				} else {
					_, _ = io.Copy(body, resp.Body)
					_ = resp.Body.Close()
				}
			}
		}
		body = nil
		return true
	}
	return false
}
