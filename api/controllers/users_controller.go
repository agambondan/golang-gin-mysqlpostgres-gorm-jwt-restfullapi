package controllers

import (
	"../auth"
	"../models"
	"github.com/gin-gonic/gin"
	"strconv"

	"net/http"
)

func (server *Server) CreateUser(c *gin.Context) {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	user := models.User{}
	err = c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	err = user.ValidateUser("")
	saveUser, err := user.SaveUser(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, saveUser)
}

func (server *Server) GetAllUser(c *gin.Context) {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	user := models.User{}
	users, err := user.FindAllUser(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Get All Data User",
		"data":    users,
	})
	return
}

func (server *Server) GetUserById(c *gin.Context) {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	id := c.Params.ByName("id")
	//uId := c.Param("id")
	uId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserById(server.DB, uint32(uId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"message": `Data By Id ` + strconv.Itoa(uId) + ` Is Found`,
		"data":    userGotten,
	})
}

func (server *Server) UpdateUserById(c *gin.Context) {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	id := c.Params.ByName("id")
	user := models.User{}
	uId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	if tokenID != uint32(uId) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "token by " + strconv.Itoa(int(tokenID)) + " not the same with " + strconv.Itoa(uId)})
		return
	}
	err = user.ValidateUser("update")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	updatedUser, err := user.UpdateUserById(server.DB, uint32(uId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data updated": updatedUser})
}

func (server *Server) DeleteUserById(c *gin.Context) {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	id := c.Params.ByName("id")
	user := models.User{}
	uId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	if tokenID != uint32(uId) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user id by " + strconv.Itoa(int(tokenID)) + " not the same with get id " + strconv.Itoa(uId)})
		return
	}
	if tokenID != 0 && tokenID != uint32(uId) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user id by " + strconv.Itoa(int(tokenID)) + " not the same with get id " + strconv.Itoa(uId)})
		return
	}
	_, err = user.DeleteUserById(server.DB, uint32(uId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "user by id" + strconv.Itoa(uId) + " is deleted"})
}
