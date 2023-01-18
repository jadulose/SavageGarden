package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("你好，世界！")
	conf, err := ReadConfig("tmp/conf.toml")
	PrintAndExit(err)
	db, err := conf.Database.Open()
	PrintAndExit(err)
	defer db.Close()
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	PrintAndExit(err)
	fmt.Println(version)
}

func PrintAndExit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
