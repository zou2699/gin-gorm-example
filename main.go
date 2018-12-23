package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
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
	//db, err = gorm.Open("mysql", "zouhl:passw0rd@tcp(192.168.3.149:3306)/test?charset=utf8&parseTime=True&loc=Local")

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

	// web api
	router.LoadHTMLGlob("web/*")
	v1 := router.Group("/web/")
	{
		v1.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{"title": "Main Website"})
		})
	}

	router.Run(":8080")
}

func DeletePerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person Person
	affectedRow := db.Where("id = ?", id).Delete(&person).RowsAffected
	if affectedRow == 0 {
		c.AbortWithStatus(404)
		c.JSON(404, gin.H{"result": "not found"})
		return
	}

	c.JSON(200, gin.H{"deletedId": id, "affectedRow": affectedRow})
}

func UpdatePerson(c *gin.Context) {
	var person Person
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "result": "not found"})
		return
	}

	c.BindJSON(&person)
	db.Save(&person)
	c.JSON(200, person)
}

func CreatePerson(c *gin.Context) {
	var person Person
	if err := c.BindJSON(&person); err != nil {
		fmt.Println(err)
	}
	if err := db.Create(&person).Error; err != nil {
		log.Println("create person err:", err)
		c.JSON(500, gin.H{"code": 500, "result": err.Error()})
		return
	}
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
