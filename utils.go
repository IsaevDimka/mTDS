package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
	"github.com/labstack/echo"
)

const utilsModuleName = "utils.go"

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func printError(header string, message interface{}, module string) {
	fmt.Fprintf(color.Output, "[ %s ]", color.RedString(header))
	fmt.Println(" ", message, " - ", module)
}

func printInfo(header string, message interface{}, module string) {
	fmt.Fprintf(color.Output, "[ %s ]", color.CyanString(header))
	fmt.Println(" ", message, " - ", module)
}

func printSuccess(header string, message interface{}, module string) {
	fmt.Fprintf(color.Output, "[ %s ]", color.GreenString(header))
	fmt.Println(" ", message, " - ", module)
}

func printDebug(header string, message interface{}, module string) {
	fmt.Fprintf(color.Output, "[ %s ]", color.YellowString(header))
	fmt.Println(" ", message, " - ", module)
}

// URIByMap Заполняем наш мап параметрами из УРИ
func URIByMap(c echo.Context, keyMap []string) map[string][]string {
	resultMap := make(map[string][]string)
	for _, item := range keyMap {
		tmp := c.Param(item)
		if tmp == "" {
			tmp = c.QueryParam(item)
		}
		resultMap[item] = append(resultMap[item], tmp)
	}
	return resultMap
}

func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

func JSONPretty(Data interface{}) string {
	var out bytes.Buffer //буфер конвертации джейсона в красивый джейсон
	jsonData, _ := json.Marshal(Data)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)
	json.Indent(&out, jsonData, "", "    ")
	return out.String()
}

// StringWithCharset this is very good function
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return StringWithCharset(length, charset)
}
