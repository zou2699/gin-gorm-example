package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

var db *gorm.DB

type Person struct {
	ID        uint   `form:"id" json:"id" xml:"id"`
	FirstName string `form:"first_name" json:"first_name" xml:"first_name" binding:"required"`
	LastName  string `form:"last_name" json:"last_name" xml:"last_name" binding:"required"`
	City      string `form:"city" json:"city" xml:"city" binding:"required"`
}

func main() {
	// = not :=
	var err error
	//db, err = gorm.Open("sqlite3", "./gorm.db")
	db, err = gorm.Open("mysql", "zouhl:passw0rd@tcp(mysql-mysql.mysql:3306)/blog?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.AutoMigrate(&Person{})

	//gin web
	log.Println(gin.Version)
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
	router.LoadHTMLGlob("templates/**/*")
	v1 := router.Group("/web/")
	{
		// get web index
		v1.GET("/users/index", func(c *gin.Context) {
			c.HTML(http.StatusOK, "users/index.tmpl", gin.H{"title": "Main Website"})
		})
		// post web index
		v1.GET("posts/index", func(c *gin.Context) {
			c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{"titile": "Add people"})
		})

	}

	_ = router.Run(":8080")
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

	_ = c.BindJSON(&person)
	db.Save(&person)
	c.JSON(200, person)
}

func CreatePerson(c *gin.Context) {
	var person Person

	if err := c.ShouldBind(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(person.FirstName)
	if err := db.Create(&person).Error; err != nil {
		log.Println("create person err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "result": err.Error()})
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
