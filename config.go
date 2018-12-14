package main

import (
	"fmt"
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
}

var cfg Config            // Конфиг инстанс
var redisdb *redis.Client // Редис

// Загрузка конфига и обработка параметров
func initConfig() {
	var actionArg string
	_ = actionArg

	err := gcfg.FatalOnly(gcfg.ReadFileInto(&cfg, configFileName))

	if len(os.Args) > 1 {
		actionArg = os.Args[1]
		if cfg.Debug.Level > 0 {
			printDebug("Current command", "actionArg", configModuleName)
		}
	}

	if err != nil {
		printError("Config error", err, configModuleName)
		os.Exit(3) // exit anyway
	} else {

		switch actionArg {

		case "config":
			{
				// минимум мы должны знать где будем слушать
				if cfg.General.Name != "" || cfg.General.Port != 0 {
					fmt.Println("[ General ]")
					if cfg.General.Name != "" {
						fmt.Println("[ -- Name ]", cfg.General.Name)
					}
					if cfg.General.Host != "" {
						fmt.Println("[ -- Host ]", cfg.General.Host)
					}
					if cfg.General.Port != 0 {
						fmt.Println("[ -- Port ]", cfg.General.Port)
					} else {
						fmt.Println("[ -- Empty ]")
					}
				}

				// минимум мы должны знать где будем слушать
				if cfg.Click.Length != 0 {
					fmt.Println("[ Click ]")
					if cfg.Click.Length != 0 {
						fmt.Println("[ -- Length ]", cfg.Click.Length)
					} else {
						fmt.Println("[ -- Empty ]")
					}
				}

				//Конфигурация редиски
				if cfg.Redis.Host != "" || cfg.Redis.Port != 0 {
					fmt.Println("[ Redis ]")
					if cfg.Redis.Host != "" {
						fmt.Println("[ -- Host ]", cfg.Redis.Host)
					}
					if cfg.Redis.Port != 0 {
						fmt.Println("[ -- Port ]", cfg.Redis.Port)
					}
					if cfg.Redis.Login != "" {
						fmt.Println("[ -- Login ]", cfg.Redis.Login)
					}
					if cfg.Redis.Password != "" {
						fmt.Println("[ -- Password ]", cfg.Redis.Password)
					} else {
						fmt.Println("[ -- Empty ]")
					}
				}

				// тут нам надо понять на каком уровне дебага мы хотим работать
				if cfg.Debug.Level != 0 {
					fmt.Println("[ Debug ]")
					if cfg.Debug.Level != 0 {
						fmt.Println("[ -- Debug level ]", cfg.Debug.Level)
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
