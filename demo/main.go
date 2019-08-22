package main

import (
	"fmt"
	"nullgo"
)

func main() {
	router := nullgo.Default()
	router.GET("/hello", func(c *nullgo.Context) {
		name := c.Query("name")
		age := c.Query("age")
		c.String("hello," + name)
		c.String("age,"+age)
	})

	router.GET("/user/:name", func(c *nullgo.Context) {
		name := c.Param("name")
		c.String("hello, " + name)
	})

	router.Websocket("/ws", nullgo.WebSocketConfig{
		OnOpen: func(context *nullgo.WebSocketContext) {
			fmt.Println("open!")
		},
		OnClose: func(context *nullgo.WebSocketContext) {
			fmt.Println("close!")
		},
		OnMessage: func(context *nullgo.WebSocketContext, s string) {
			fmt.Println(context.Conn.RemoteAddr(), "send message:", s)
		},
		OnError: nil,
	})


	router.POST("/abc", func(c *nullgo.Context) {
		message := c.PostV("message")
		c.String("message:" + message)
	})




	router.Run(":8080")

	//GetInfo()
}










