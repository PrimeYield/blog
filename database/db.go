package database

import (
	"context"
	"fmt"
	"log"
	"practise/models"
	"practise/pkg/setting"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// db:=Client.Database(databaseName)
// collection:=db.Collection(collectionName)
func GetCollection(collectionName string) *mongo.Collection {
	if Client == nil {
		log.Fatal("MongoDB client is not initialized. Call MongodbJoin first.")
	}
	return Client.Database(databaseName).Collection(collectionName)
}

func CreateUser(username string, age int) (primitive.ObjectID, error) {
	collection := GetCollection("users")
	user := models.User{
		Username: username,
		Age:      age,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// InsertOne 會插入單個文檔
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert user: %v", err)
	}
	// 返回插入文檔的 ID
	return result.InsertedID.(primitive.ObjectID), nil
}

// Read by ID
func FindUserByID(id primitive.ObjectID) (*models.User, error) {
	collection := GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	//"_id" 是bson tag的命名
	filter := bson.M{"_id": id} //M is an unordered representation of a BSON document.
	//FindOne executes a find command and returns a SingleResult for one document in the collection.
	err := collection.FindOne(ctx, filter).Decode(&user) //Decode will unmarshal the document represented by this SingleResult into v(&user).
	if err != nil {
		//ErrNoDocuments is returned by SingleResult methods when the operation that created the SingleResult did not return any documents.
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by ID: %v", err)
	}
	return &user, nil
}

// find multipy users by username
func FindUsersByUsername(username string) ([]models.User, error) {
	collection := GetCollection("users") //這個是mongodb的起手式
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	//Find executes a find command and returns a Cursor over the matching documents in the collection.
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find users by username: %v", err)
	}
	defer cursor.Close(ctx)
	//存儲Find結果
	var users []models.User
	//All iterates the cursor and decodes each document into results(&users).
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %v", err)
	}
	return users, nil
}

func UpdateUserByID(id primitive.ObjectID, updates bson.M) (int64, error) {
	collection := GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	// $set 操作符用於設定或替換欄位值
	update := bson.M{"$set": updates}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("failed to update user: %v", err)
	}
	return result.ModifiedCount, nil // 返回被修改的文檔數量
}

func DeleteUserByID(id primitive.ObjectID) (int64, error) {
	collection := GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete user: %v", err)
	}
	return result.DeletedCount, nil // 返回被刪除的文檔數量
}

// Articles
// Get
func GetArticles(CreatedBy string) ([]models.Article, error) {
	collection := GetCollection("articles")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"created_by": CreatedBy}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find articles by username:%v", err)
	}
	defer cursor.Close(ctx)

	var articles []models.Article
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, fmt.Errorf("failed to decode articles:%v", err)
	}
	return articles, nil
}

func GetArticle(id primitive.ObjectID) (*models.Article, error) {
	collection := GetCollection("articles") //其實這個也可以用viper管理，之後再優化
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var article models.Article
	filter := bson.M{"_id": id}
	err := collection.FindOne(ctx, filter).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find article :%v", err)
	}
	return &article, nil
}

// Created
func CreateArticle(article models.Article) (primitive.ObjectID, error) {
	collection := GetCollection("articles")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	article.CreatedAt = time.Now()

	result, err := collection.InsertOne(ctx, article)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert article: %w", err)
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// Update
func UpdateArticle(id primitive.ObjectID, updateData bson.M) (int64, error) {
	collection := GetCollection("articles")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	finalUpdate := bson.M{
		"$set": updateData, // $set 來自請求體提供的數據
	}

	if _, ok := updateData["$set"]; ok {
		finalUpdate = updateData
	} else {
		finalUpdate = bson.M{
			"$set": updateData,
		}
	}

	if setMap, ok := finalUpdate["$set"].(bson.M); ok {
		setMap["updated_at"] = time.Now()
	} else {
		finalUpdate["$set"] = bson.M{"updated_at": time.Now()}
	}

	result, err := collection.UpdateOne(ctx, filter, finalUpdate)
	if err != nil {
		return 0, fmt.Errorf("failed to update article: %v", err)
	}
	return result.ModifiedCount, nil
}

// Delete
func DelArticle(id primitive.ObjectID) (bool, error) {
	collection := GetCollection("articles")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to delete article: %v", err)
	}
	return true, nil
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
