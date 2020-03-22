package controllers

import (
	"../auth"
	"../models"
	"github.com/gin-gonic/gin"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err = user.ValidateUser("login")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	token, err := server.SignIn(user.Email, user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Username or Password is wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (server *Server) SignIn(email, username, password string) (string, error) {
	user := models.User{}
	err := server.DB.Debug().Model(models.User{}).Where("email = ? OR username = ?", email, username).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(uint32(user.ID))
}
