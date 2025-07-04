package handlers

import (
	"fmt"
	"log"
	"net/http"
	"practise/database"
	"practise/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateArticleHandler(c *gin.Context) {
	fmt.Println("CreatHandler", c.Keys)
	// fmt.Printf("CreatHandler %v,%b", c.Get("username"))
	fmt.Println("CreatHandler", c.GetString("username"))

	var article models.Article

	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %v", err.Error())})
		return
	}

	article.CreatedBy = c.GetString("username")

	// username, exists := c.Get("username")
	// if !exists {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Authenticated user information not found",
	// 	})
	// 	return
	// }
	// models.UserArticle.CreatedBy = username.(string)

	// 完成一個認證系統

	// article.CreatedBy = c.GetString("username")
	// fmt.Println(c.GetString("JwtID"))
	// fmt.Println(c.Keys)
	// article.CreatedBy = c.

	insertedID, err := database.CreateArticle(article)
	if err != nil {
		log.Printf("Failed to create article: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create article in database",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Article created successfully",
		"id":         insertedID.Hex(),
		"title":      article.Title,
		"content":    article.Content,
		"created_by": article.CreatedBy,
		"created_at": article.CreatedAt,
	})
}

func UpdateArticleHandler(c *gin.Context) {
	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid article ID format",
		})
		return
	}

	var updates bson.M
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %v", err.Error())})
		return
	}

	// 確保不允許前端直接修改 CreatedAt、UpdatedAt 或 _id 等字段
	delete(updates, "created_at") // 移除前端可能傳來的 created_at
	delete(updates, "updated_at")
	delete(updates, "id")
	delete(updates, "CreatedBy")

	// 檢查有沒有需要更新的字段 (==0時，只有 updated_at 會自動設置)
	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No specific fields provided for update, only UpdatedAt will be set.",
			"id":      objID.Hex(),
		})
		return
	}

	// db update
	modifiedCount, err := database.UpdateArticle(objID, updates)
	if err != nil {
		log.Printf("Failed to update article: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update article in database",
		})
		return
	}

	if modifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Article not found or no changes made (other than UpdatedAt)",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":       "Article updated successfully",
			"id":            objID.Hex(),
			"modifiedCount": modifiedCount,
		})
	}
}

func GetAuthorArticlesHandler(c *gin.Context) {
	createdBy := c.Param("createdBy") //author似乎是更好的單字XD

	res, err := database.GetArticles(createdBy)
	if err != nil {
		log.Printf("Failed to get articles:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get articles from database",
		})
		return
	}
	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Author not exist",
		})
		return
	}
	c.JSON(http.StatusAccepted, res)
}

func GetArticleHandler(c *gin.Context) {
	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid article ID format",
		})
		return
	}
	res, err := database.GetArticle(objID)
	if err != nil {
		log.Printf("Failed to get article: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get article from database",
		})
		return
	}
	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Article not found",
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func DeleteArticleHandler(c *gin.Context) {
	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID format"})
		return
	}

	deletedCount, err := database.DelArticle(objID)
	if err != nil {
		log.Printf("Failed to delete article: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article from database"})
		return
	}

	if !deletedCount {
		c.JSON(http.StatusNotFound, gin.H{"message": "Article not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully", "id": objID.Hex()})
	}
}
