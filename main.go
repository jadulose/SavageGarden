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
	err = db.Ping()
	PrintAndExit(err)

	mail, err := conf.Mail.CreateTestMsg("peilin.fan@fansionia.xyz")
	PrintAndExit(err)
	mClient, err := conf.Mail.Open()
	PrintAndExit(err)
	fmt.Println("正在发送邮件...")
	err = mClient.DialAndSend(mail)
	PrintAndExit(err)
	fmt.Println("已成功发送邮件")
}

func PrintAndExit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
