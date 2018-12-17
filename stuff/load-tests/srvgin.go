package main

import ("fmt";"github.com/gin-gonic/gin")

var counter int

func main() {
	counter = 0
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		counter++;
		fmt.Println(counter)

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":80") // listen and serve on 0.0.0.0:8080
}