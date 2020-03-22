package controllers

import (
	"../auth"
	"../models"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (server *Server) CreatePost(c *gin.Context) {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	post := models.Post{}
	err = post.Validate()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	if post.AuthorID != uid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	savePost, err := post.SavePost(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": savePost})
}

func (server *Server) GetAllPost(c *gin.Context) {
	post := models.Post{}
	findAllPost, err := post.FindAllPost(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": findAllPost})
}

func (server *Server) GetPostById(c *gin.Context) {
	id := c.Params.ByName("id")
	pId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	post := models.Post{}
	findPost, err := post.FindPostByID(server.DB, uint64(pId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusFound, gin.H{"data": findPost})
}

func (server *Server) UpdatePostById(c *gin.Context) {
	id := c.Params.ByName("id")
	// Check if the post id is valid
	pId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	//CHeck if the auth token is valid and get the user id from it
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	// Check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pId).Take(&post).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	// If a user attempt to update a post not belonging to him
	if uid != post.AuthorID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": errors.New("UnAuthorized").Error()})
		return
	}
	// Read the data posted
	err = c.BindJSON(post)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	// Start processing the request data
	postUpdate := models.Post{}
	//Also check if the request user id is equal to the one gotten from token
	if uid != postUpdate.AuthorID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": errors.New("UnAuthorized").Error()})
		return
	}
	err = postUpdate.Validate()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	postUpdate.ID = post.ID //this is important to tell the model the post id to update, the other update field are set above
	postUpdated, err := postUpdate.UpdatePostById(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": postUpdated})
}

func (server *Server) DeletePostById(c *gin.Context) {
	id := c.Params.ByName("id")
	pId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	uId, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	post := models.Post{}
	err = server.DB.Debug().Model(&models.Post{}).Where("id = ?", pId).Take(&post).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if uId != post.AuthorID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "get Id and params not the same"})
		return
	}
	_, err = post.DeletePostById(server.DB, uint64(pId), uId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "data by id " + strconv.Itoa(pId) + " success deleted"})
}
