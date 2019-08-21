package main

import (
	"fmt"
	. "nullgo"
)

func GetInfo() {
	cfg, err := LoadConfig("test.conf")
	if err != nil {
		Error("get info error")
		return
	}
	name := cfg.String("username")
	pwd := cfg.String("password")
	DBName := cfg.String("DBName")

	fmt.Println("name:", name)
	fmt.Println("password:", pwd)
	fmt.Println("DataBaseName:", DBName)
}
