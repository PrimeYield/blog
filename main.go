package main

import (
	"context"
	"fmt"
	"practise/database"
	"practise/global"
	"practise/pkg/setting"
)

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		panic(err)
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	setupSetting()
	err := database.MongodbJoin(&global.DatabaseSetting)
	if err != nil {
		panic("Connect to db failed:" + err.Error())
	}
	defer func() {
		err := database.Client.Disconnect(context.Background())
		if err != nil {
			//todo
		}
	}()
}
