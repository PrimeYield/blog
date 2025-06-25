package main

import (
	"context"
	"fmt"
	"practise/database"
	"practise/global"
	"practise/handlers"
	"practise/pkg/setting"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	port := global.ServerSetting.Port
	dbGroup := r.Group("/db")
	{
		dbGroup.POST("/add", handlers.CreateUserHandler) //包裝db的func
		dbGroup.GET("/getuser/:id", handlers.GetUserHandler)
		dbGroup.GET("/updateuser/:id", handlers.UpdateUserHandler)
		dbGroup.DELETE("/deluser/:id", handlers.DelUserHandler)
	}

	r.Run(":" + port)
}
