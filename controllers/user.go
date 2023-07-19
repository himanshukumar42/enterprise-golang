package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/himanshukumar42/enterprise/forms"
	"github.com/himanshukumar42/enterprise/models"
)

// UserController...
type UserController struct{}

var userModel = new(models.UserModel)
var userForm = new(forms.UserForm)

func getUserID(c *gin.Context) (userID int64) {
	// MustGet returns the value for the given key if it exists, otherwise it panics
	return c.MustGet("userID").(int64)
}

func (ctrl UserController) Login(c *gin.Context) {
	var loginForm forms.LoginForm

	if validationErr := c.ShouldBindJSON(&loginForm); validationErr != nil {
		message := userForm.Login(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	user, token, err := userModel.Login(loginForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "invalid login details"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "user": user, "token": token})
}

func (ctrl UserController) Register(c *gin.Context) {
	var registerForm forms.RegisterForm

	if validationErr := c.ShouldBindJSON(&registerForm); validationErr != nil {
		message := userForm.Login(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	user, err := userModel.Register(registerForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Registered", "user": user})
}

func (ctrl UserController) Logout(c *gin.Context) {
	au, err := authModel.ExtractTokenMetaData(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not logged in"})
		return
	}

	deleted, delErr := authModel.DeleteAuth(au.AccessUUID)
	if delErr != nil || deleted == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})

}
