package database

import (
	"context"
	"fmt"
	"log"
	"practise/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// db:=Client.Database(databaseName)
// collection:=db.Collection(collectionName)
func GetCollection(collectionName string) *mongo.Collection {
	if Client == nil {
		log.Fatal("MongoDB client is not initialized. Call MongodbJoin first.")
	}
	return Client.Database(databaseName).Collection(collectionName)
}

func CreateUser(username, password string, age int) (primitive.ObjectID, error) {
	if len(password) < 6 {
		return primitive.NilObjectID, fmt.Errorf("len(password) > 6")
	}

	target, _ := FindUserByUsername(username)
	if target != nil { //如果username是FindUserByUsername無法匹配的，會回傳nil,error才對 不明白為什麼target會有值
		return primitive.NilObjectID, fmt.Errorf("%v is already exists", username)
	}
	// if target.Username != "" {
	// 	return primitive.NilObjectID, fmt.Errorf("%s is already exists",target.Username)
	// }
	// if err != nil && err.Error() == fmt.Errorf("username: '%s'  is not exists", username).Error() {
	// 	//do nothing
	// } else if err != nil {
	//     return primitive.NilObjectID, fmt.Errorf("failed to check existing user: %w", err)
	// } else if target != nil {
	//     return primitive.NilObjectID, fmt.Errorf("%s is already exists", username)
	// }
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to hash password: %v", err)
	}
	collection := GetCollection("users")
	// fmt.Println(collection)
	user := models.User{
		Username:  username,
		Password:  string(hashedPassword),
		Age:       age,
		CreatedAt: time.Now(),
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
	fmt.Println(user)
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
func FindUserByUsername(username string) (*models.User, error) {
	collection := GetCollection("users") //這個是mongodb的起手式
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	//Find executes a find command and returns a Cursor over the matching documents in the collection.
	cursor := collection.FindOne(ctx, filter)

	var user models.User
	err := cursor.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments { //這裡return 了 nil 怎麼會在Create()裡面有值呢???
			return nil, fmt.Errorf("username: '%s'  is not exists", username)
		}
		return nil, fmt.Errorf("bad request for %s", username)
	}
	return &user, nil
}

func UpdateUserByID(id primitive.ObjectID, updates bson.M) (int64, error) {
	collection := GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	// fmt.Printf("filter: %v\n", filter)
	// $set 操作符用於設定或替換欄位值
	// updateAt := time.Now()
	// updateAt := time.Now()
	// update := bson.M{"$set": updates, "updated_at": updateAt}
	setUpdate := make(bson.M)
	for k, v := range updates {
		setUpdate[k] = v
	}
	updateAt := time.Now()
	setUpdate["updated_at"] = updateAt
	newPassword, err := bcrypt.GenerateFromPassword([]byte(setUpdate["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %v", err)
	}
	setUpdate["password"] = newPassword
	update := bson.M{"$set": setUpdate}

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
