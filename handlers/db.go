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

func CreateUserHandler(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	insertedID, err := database.CreateUser(user.Username, user.Age)
	if err != nil {
		log.Printf("Failed to create user:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user in database",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":  "User created successfully",
		"id":       insertedID.Hex(), //Hex returns the hex encoding of the ObjectID as a string.
		"username": user.Username,
		"age":      user.Age,
	})
}

func GetUserHandler(c *gin.Context) {
	idStr := c.Param("id") //router: /db/getuser/:id

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	res, err := database.FindUserByID(objID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user from database",
		})
		return
	}
	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not exist",
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"id":       res.ID.Hex(), //Hex returns the hex encoding of the ObjectID as a string.
		"username": res.Username,
		"age":      res.Age,
	})
}

func UpdateUserHandler(c *gin.Context) {
	idStr := c.Param("id") //router: //db/updateuser/:id
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}
	var updates bson.M
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %v", err.Error()),
		})
		return
	}
	modifiedCount, err := database.UpdateUserByID(objID, updates)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user in database",
		})
		return
	}
	if modifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found or no changes made",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":       "User updated successfully",
			"id":            objID.Hex(),
			"modifiedCount": modifiedCount,
		})
	}
}

func DelUserHandler(c *gin.Context) {
	idStr := c.Param("id") //router: /db/deluser/:id
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}
	modifiedCount, err := database.DeleteUserByID(objID)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user in database",
		})
		return
	}
	if modifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "User delete successfully",
		"id":            objID.Hex(),
		"modifiedCount": modifiedCount,
	})
}
