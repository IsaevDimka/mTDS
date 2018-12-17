package main

import (
	"net/http"
	"fmt"
	"github.com/labstack/echo"
	_"github.com/labstack/echo/middleware"
)

var counter int

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	e.Start(":80")
	// Start server
	//e.Logger.Fatal()
}

// Handler
func hello(c echo.Context) error {

	for i:=0; i<100000; i++ {
		counter++;
		fmt.Println(counter)
	}


/*	go func() {
	//	c.String(http.StatusOK, "Hello, World!")
		c.Redirect(302, "https://lenta.ru/")
	}()
*/

//	var e error
		
	return c.String(http.StatusOK, "Hello, World!") //fmt.Errorf("math: square root of negative number")
}