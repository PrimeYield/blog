package database

import (
	"context"
	"fmt"
	"practise/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Articles
// Get
// 可以考慮新增一個GetArticlesByTitle
// 要思考如何將密碼藏起來不曝露
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

// Created  不在乎title是否重覆
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

	//修正無法正確寫入updated_at
	updateAt := time.Now()
	includeUpdateAt := make(bson.M)
	for k,v := range updateData {
		includeUpdateAt[k] = v
	}
	includeUpdateAt["updated_at"] = updateAt
	
	finalUpdate := bson.M{
		"$set": includeUpdateAt, // $set 來自請求體提供的數據
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