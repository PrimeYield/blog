package database

import (
	"context"
	"log"
	"practise/pkg/setting"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var databaseName string

func MongodbJoin(databaseSetting *setting.DatabaseSetting) error {
	//connect DB
	clientOptions := options.Client().ApplyURI("mongodb://" + databaseSetting.MongodbHost + ":" + databaseSetting.MongodbPort)

	//使用連接選項連接
	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	databaseName = databaseSetting.Mongodb_db // 儲存資料庫名稱以供後續使用
	log.Println("Connected to MongoDB!")
	return nil
}

// func setupSetting() error {
// 	setting, err := setting.NewSetting()
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = setting.ReadSection("Server", &global.ServerSetting)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	err = setting.ReadSection("Database", &global.DatabaseSetting)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	return nil
// }

// func main() {
// 	err := setupSetting()
// 	if err != nil {
// 		panic(err)
// 	}
// 	MongodbJoin(&global.DatabaseSetting)
// 	defer Client.Disconnect(context.Background())
// 	// CreateUser("pig", 100) ok
// 	// idString := "685aa7d96195d8c36968bc6f"
// 	// objID, err := primitive.ObjectIDFromHex(idString)
// 	// if err != nil {
// 	// 	log.Fatalf("Invalid ObjectID string: %v", err)
// 	// }
// 	// fmt.Println(FindUserByID(objID)) //&{ObjectID("685aa7d96195d8c36968bc6f") pig 100} <nil>
// 	// fmt.Println(FindUsersByUsername("pig")) //[{ObjectID("685aa7d96195d8c36968bc6f") pig 100} {ObjectID("685aa7ef4ac2ba359384140a") pig 100}] <nil>
// 	// idString := "685aa7d96195d8c36968bc6f"
// 	// objID, err := primitive.ObjectIDFromHex(idString)
// 	// if err != nil {
// 	// 	log.Fatalf("Invalid ObjectID string: %v", err)
// 	// }
// 	// fmt.Println(DeleteUserByID(objID)) //1 <nil>
// }
