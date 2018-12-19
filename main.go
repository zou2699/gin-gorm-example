package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var db *gorm.DB

type Person struct {
	ID       uint   `json:"id"`
	FistName string `json:"fist_name"`
	LastName string `json:"last_name"`
	City     string `json:"city"`
}

func main() {
	// = not :=
	var err error
	db, err = gorm.Open("sqlite3", "./gorm.db")
	//db, err = gorm.Open("mysql", "zouhl:passw0rd@tcp(192.168.3.149:3306)/blog?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.AutoMigrate(&Person{})

	router := gin.Default()
	//ping for health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/people", GetPeople)
	router.GET("/people/:id", GetPerson)
	router.POST("/people", CreatePerson)
	router.PUT("/people/:id", UpdatePerson)
	router.DELETE("/people/:id", DeletePerson)

	router.Run(":8080")
}

func DeletePerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person Person
	if err := db.Where("id = ?", id).Delete(&person).Error;err!=nil {
		c.AbortWithStatus(404)
		log.Println(err)
		return
	}
	c.JSON(200, gin.H{"id #" + id: "delete"})
}

func UpdatePerson(c *gin.Context) {
	var person Person
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		log.Println(err)
	}
	c.BindJSON(&person)
	db.Save(&person)
	c.JSON(200, person)
}

func CreatePerson(c *gin.Context) {
	var person Person
	if err := c.BindJSON(&person); err!=nil {
		fmt.Println(err)
	}
	db.Create(&person)
	c.JSON(200, person)
}

func GetPerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person Person
	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		log.Println(err)
	} else {
		c.JSON(200, person)
	}
}

func GetPeople(c *gin.Context) {
	var people []Person
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		log.Println(err)
	} else {
		c.JSON(200, people)
	}
}
