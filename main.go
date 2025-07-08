package main

import (
	"context"
	"fmt"
	"practise/database"
	"practise/global"
	"practise/handlers"
	"practise/middleware"
	"practise/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwa"
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

	err = setting.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		fmt.Println(err)
		return err
	}
	global.JWTSetting.Algorithm = jwa.HS256
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

	// err = config.LoadJWTSigningKey(global.JWTSetting.Secret)
	// if err != nil {
	// 	log.Fatalf("Failed to load JWT signing key: %v", err)
	// }

	r := gin.Default()
	r.Static("public", "./public")
	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})
	port := global.ServerSetting.Port
	r.POST("/login", handlers.LoginHandler)
	userGroup := r.Group("/user")
	{
		userGroup.POST("/add", handlers.CreateUserHandler) //包裝db的func
		userGroup.GET("/get/:id", handlers.GetUserHandler)
		userGroup.POST("/update/:id", handlers.UpdateUserHandler)
		userGroup.DELETE("/delete/:id", handlers.DelUserHandler)
	}
	articleGroup := r.Group("/arti")
	articleGroup.Use(middleware.AuthMiddleware())
	{
		articleGroup.POST("/create", handlers.CreateArticleHandler)
		articleGroup.POST("/update", handlers.UpdateArticleHandler)
		articleGroup.DELETE("/delete/:id", handlers.DeleteArticleHandler)
		articleGroup.POST("/getById/:id", handlers.GetArticleHandler)
		articleGroup.POST("/getByAuthor", handlers.GetAuthorArticlesHandler)
	}

	r.Run(":" + port)
}
