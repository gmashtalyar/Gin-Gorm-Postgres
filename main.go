package main

import (
	"Gin_Gorm_Postgres/models"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDB() {
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	DB.AutoMigrate(&models.User{})
}

func main() {
	router := gin.Default()
	initDB()
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/", func(c *gin.Context) {
		var users []models.User
		DB.Find(&users)
		c.HTML(http.StatusOK, "index.html", gin.H{"users": users})
	})

	router.GET("/users/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.html", nil)
	})

	router.POST("/users", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBind(&user); err != nil { // ShouldBindJSON
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if result := DB.Create(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
		// c.Redirect(http.StatusMovedPermanently, "/")

	})

	router.GET("/users/edit/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := DB.First(&user, id).Error; err != nil {
			c.String(http.StatusNotFound, "User not found")
			return
		}
		c.HTML(http.StatusOK, "edit.html", gin.H{"user": user})
	})

	router.POST("/users/update/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := DB.First(&user, id).Error; err != nil {
			c.String(http.StatusNotFound, "User not found")
			return
		}
		if err := c.ShouldBind(&user); err != nil {
			c.HTML(http.StatusBadRequest, "edit.html", gin.H{"error": err.Error(), "user": user})
			return
		}

		DB.Save(&user)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	router.GET("/users/delete/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := DB.Delete(&user, id).Error; err != nil {
			c.String(http.StatusNotFound, "User not found")
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	router.Run(":8080")
}
