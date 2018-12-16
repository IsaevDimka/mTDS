package config

import (
	"fmt"
	"metatds/utils"
	"os"

	"github.com/go-redis/redis"
	gcfg "gopkg.in/gcfg.v1"
)

const configFileName = "settings.ini"
const configModuleName = "config.go"

// Config Тип для хранения конфига
type Config struct {
	General struct {
		Name string
		Host string
		Port int
	}
	Click struct {
		Length int
	}
	Redis struct {
		Host     string
		Port     int
		Login    string
		Password string
	}
	Debug struct {
		Level int
	}
	Telegram struct {
		Socks5Proxy    string
		Socks5User     string
		Socks5Password string
		ApiURL         string
		Token          string
		Recipients     string
	}
}

var Cfg Config                     // Конфиг инстанс
var Redisdb *redis.Client          // Редис
var Telegram utils.TelegramAdapter // инстанс бота

// Загрузка конфига и обработка параметров
func InitConfig() {
	var actionArg string
	_ = actionArg

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
					} else {
						fmt.Println("[ -- Empty ]")
					}
				}

				// минимум мы должны знать где будем слушать
				if Cfg.Click.Length != 0 {
					fmt.Println("[ Click ]")
					if Cfg.Click.Length != 0 {
						fmt.Println("[ -- Length ]", Cfg.Click.Length)
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
					} else {
						fmt.Println("[ -- Empty ]")
					}
				} else {
					fmt.Println("[ Empty configuration ]")
				}

				// тут нам надо понять на каком уровне дебага мы хотим работать
				if Cfg.Telegram.Token != "" && Cfg.Telegram.Socks5Proxy != "" {
					fmt.Println("[ Telegram ]")
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
