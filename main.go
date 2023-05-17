package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type User struct {
	gorm.Model
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open(postgres.Open("postgres://postgres:4451122@localhost:5432/golang-db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{})

	r := gin.Default()

	r.POST("/register", CreateUser)
	r.POST("/login", LoginUser)

	r.Run("localhost:8080")
}

func CreateUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	db.Create(&user)
	c.JSON(http.StatusCreated, user)
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginUser(c *gin.Context) {
	var userLogin UserLogin
	var foundUser User

	if err := c.BindJSON(&userLogin); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	fmt.Println(userLogin)

	if result := db.Where("email = ? AND password = ?", userLogin.Email, userLogin.Password).First(&foundUser); result.Error != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, foundUser)
}
