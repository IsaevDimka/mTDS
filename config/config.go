/****************************************************************************************************
*
* Config init module, special for Meta CPA, Ltd.
* by Michael S. Merzlyakov AFKA predator_pc@06122018
* version v2.0.3
*
* created at 06122018
* last edit: 16122018
*
*****************************************************************************************************/

package config

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis"
	"gopkg.in/gcfg.v1"
	"metatds/utils"
	"os"
	"time"
)

const configFileName = "settings.ini"
const configModuleName = "config.go"

// Config Тип для хранения конфига
type Config struct {
	General struct {
		Name       string
		Host       string
		Port       int
		ConfReload int
		OS         string
	}
	Click struct {
		Length         int
		ApiUrl         string
		ApiToken       string
		DropToRedis    int
		DropFilesToAPI int
	}
	Redis struct {
		Host     string
		Port     int
		Login    string
		Password string
	}
	Debug struct {
		Test  bool
		Level int
	}
	Telegram struct {
		MsgInterval    int
		UseProxy       bool
		Socks5Proxy    string
		Socks5User     string
		Socks5Password string
		ApiURL         string
		Token          string
		Recipients     string
	}
}

var Cfg Config            // Конфиг инстанс
var Redisdb *redis.Client // Редис
var IsRedisAlive = false
var Telegram utils.TelegramAdapter // инстанс бота
var TDSStatistic utils.TDSStats    // Инстанс статистики
var ResponseAverage []time.Duration

// Загрузка конфига и обработка параметров
func InitConfig() {
	var actionArg string
	_ = actionArg

	//fp, _ := filepath.Abs(configFileName)
	err := gcfg.FatalOnly(gcfg.ReadFileInto(&Cfg, configFileName))

	if len(os.Args) > 1 {
		actionArg = os.Args[1]
		if Cfg.Debug.Level > 0 {
			utils.PrintDebug("Current command", "actionArg", configModuleName)
		}
	}

	if err != nil {
		utils.PrintError("Config error", err, configModuleName)
		os.Exit(3) // exit anyway
	} else {

		switch actionArg {

		case "config":
			{
				// минимум мы должны знать где будем слушать
				if Cfg.General.Name != "" || Cfg.General.Port != 0 {
					fmt.Println("[ General ]")
					if Cfg.General.Name != "" {
						fmt.Println("[ -- Name ]", Cfg.General.Name)
					}
					if Cfg.General.Host != "" {
						fmt.Println("[ -- Host ]", Cfg.General.Host)
					}
					if Cfg.General.Port != 0 {
						fmt.Println("[ -- Port ]", Cfg.General.Port)
					}
					if Cfg.General.ConfReload != 0 {
						fmt.Println("[ -- Config reload interval ]", Cfg.General.ConfReload)
					}
					if Cfg.General.OS != "" {
						fmt.Println("[ -- OS ]", Cfg.General.OS)
					} else {
						fmt.Println("[ -- Empty ]")
					}
				}

				// минимум мы должны знать где будем слушать
				if Cfg.Click.Length != 0 {
					fmt.Println("[ Click ]")
					if Cfg.Click.Length != 0 {
						fmt.Println("[ -- Length ]", Cfg.Click.Length)
					}
					if Cfg.Click.ApiUrl != "" {
						fmt.Println("[ -- ApiUrl ]", Cfg.Click.ApiUrl)
					}
					if Cfg.Click.ApiToken != "" {
						fmt.Println("[ -- ApiToken ]", Cfg.Click.ApiToken)
					}
					if Cfg.Click.DropToRedis != 0 {
						fmt.Println("[ -- DropToRedis ]", Cfg.Click.DropToRedis)
					}
					if Cfg.Click.DropFilesToAPI != 0 {
						fmt.Println("[ -- DropFilesToAPI ]", Cfg.Click.DropFilesToAPI)
					} else {
						fmt.Println("[ -- Empty ]")
					}
				}

				//Конфигурация редиски
				if Cfg.Redis.Host != "" || Cfg.Redis.Port != 0 {
					fmt.Println("[ Redis ]")
					if Cfg.Redis.Host != "" {
						fmt.Println("[ -- Host ]", Cfg.Redis.Host)
					}
					if Cfg.Redis.Port != 0 {
						fmt.Println("[ -- Port ]", Cfg.Redis.Port)
					}
					if Cfg.Redis.Login != "" {
						fmt.Println("[ -- Login ]", Cfg.Redis.Login)
					}
					if Cfg.Redis.Password != "" {
						fmt.Println("[ -- Password ]", Cfg.Redis.Password)
					} else {
						fmt.Println("[ -- Empty ]")
					}
				}

				// тут нам надо понять на каком уровне дебага мы хотим работать
				if Cfg.Debug.Level != 0 {
					fmt.Println("[ Debug ]")
					if Cfg.Debug.Level != 0 {
						fmt.Println("[ -- Debug level ]", Cfg.Debug.Level)
					}
					if Cfg.Debug.Test != false {
						fmt.Println("[ -- Debug TEST mode]", Cfg.Debug.Test)
					} else {
						fmt.Println("[ -- Empty ]")
					}
				} else {
					fmt.Println("[ Empty configuration ]")
				}

				// тут нам надо понять на каком уровне дебага мы хотим работать
				if Cfg.Telegram.Token != "" && Cfg.Telegram.Socks5Proxy != "" {
					fmt.Println("[ Telegram ]")
					if Cfg.Telegram.MsgInterval != 0 {
						fmt.Println("[ -- Statistic sending interval ]", Cfg.Telegram.MsgInterval)
					}
					if Cfg.Telegram.UseProxy != false {
						fmt.Println("[ -- Use Proxy ]", Cfg.Telegram.UseProxy)
					}
					if Cfg.Telegram.Socks5Proxy != "" {
						fmt.Println("[ -- Telegram proxy ]", Cfg.Telegram.Socks5Proxy)
					}
					if Cfg.Telegram.Socks5User != "" {
						fmt.Println("[ -- Telegram proxy user ]", Cfg.Telegram.Socks5User)
					}
					if Cfg.Telegram.Socks5Password != "" {
						fmt.Println("[ -- Telegram proxy password ]", Cfg.Telegram.Socks5Password)
					}
					if Cfg.Telegram.ApiURL != "" {
						fmt.Println("[ -- Telegram Api URL]", Cfg.Telegram.ApiURL)
					}
					if Cfg.Telegram.Token != "" {
						fmt.Println("[ -- Telegram Token ]", Cfg.Telegram.Token)
					}
					if Cfg.Telegram.Recipients != "" {
						fmt.Println("[ -- Telegram Recipients ]", Cfg.Telegram.Recipients)
					} else {
						fmt.Println("[ -- Empty ]")
					}
				} else {
					fmt.Println("[ Empty configuration ]")
				}
				os.Exit(3) // exit anyway
			}
		case "run":
			break
		default:
			{
				fmt.Println("Usage: [this-file] command options")
				fmt.Println("Commands: --run - start API service")
				fmt.Println("          --config - show usable ini file settings")
				fmt.Println("          --help /none - show this message")
				os.Exit(3)
			}
		}
	}
}

func ReloadConfig() {
	err := gcfg.FatalOnly(gcfg.ReadFileInto(&Cfg, configFileName))

	if err != nil {
		if Cfg.Debug.Level > 1 {
			utils.PrintError("Reload config error", err, configModuleName)
			utils.PrintInfo("Current config", "", configModuleName)
			spew.Dump(Cfg)
		}
	} else {
		if Cfg.Debug.Level > 1 {
			utils.PrintDebug("Config reloaded", Cfg, configFileName)
			spew.Dump(Cfg)
		}
	}
}
