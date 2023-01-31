package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func M(c *gin.Context) {
	c.Next()
}

func Test(c *gin.Context) {
	fmt.Printf("%v\r\n", c.HandlerNames())
	c.String(200, "hello gin")
	// c.Abort()
}

func Test2(c *gin.Context) {
	fmt.Printf("%v\r\n", c.HandlerNames())
	c.String(400, "hello test2")
}

func main() {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.Use(M)
		v1.GET("/Test", Test, Test2)
		v1.GET("/Recover", func(c *gin.Context) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("%v\r\n", "错误恢复")
				}
			}()
			panic("错误")
		})
	}

	router.Run(":8080")
}
