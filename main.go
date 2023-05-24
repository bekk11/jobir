package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type User struct {
	gorm.Model
	FirstName string `json:"firstName" form:"firstName"`
	LastName  string `json:"lastName" form:"lastName"`
	Email     string `json:"email" form:"email"`
	Password  string `json:"password" form:"password"`
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

	r.GET("/list", ListUser)
	r.POST("/register", CreateUser)
	r.POST("/login", LoginUser)

	err := r.Run("localhost:8080")
	if err != nil {
		return
	}
}

func CreateUser(c *gin.Context) {
	var user User
	if err := c.Bind(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	db.Create(&user)
	c.JSON(http.StatusCreated, user)
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ListUser(c *gin.Context) {
	var users []User

	if result := db.Find(&users); result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
	}

	c.JSON(http.StatusOK, &users)
}

func LoginUser(c *gin.Context) {
	var userLogin UserLogin
	var foundUser User

	if err := c.BindJSON(&userLogin); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if result := db.Where("email = ? AND password = ?", userLogin.Email, userLogin.Password).First(&foundUser); result.Error != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, foundUser)
}
