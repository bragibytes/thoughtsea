package controllers

import (
	"dedpidgon/thoughtsea/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type posts struct{}

func (x posts) create(c *gin.Context) {

	var post *models.Post
	if err := c.Bind(&post); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	oid, err := primitive.ObjectIDFromHex(c.GetHeader("Authorization"))
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
	}
	post.Author = oid
	if err := post.Save(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (x posts) read(c *gin.Context) {
	list, err := models.Post{}.GetAll()
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}

func (x posts) readOne(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	post, err := models.Post{ID: oid}.Get()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, post)
}

func (x posts) update(c *gin.Context) {

	var post *models.Post
	if err := c.Bind(&post); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if !usersMatch(c, post.Author) {
		c.String(http.StatusUnauthorized, "i cant let you do that")
		return
	}
	if err := post.Update(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, post)
}

func (x posts) destroy(c *gin.Context) {

	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	post := &models.Post{ID: oid}
	if !usersMatch(c, post.Author) {
		c.String(http.StatusUnauthorized, "i cant let you do that")
		return
	}
	if err := post.Destroy(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, post)
}

func (x posts) vote(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response{false, err.Error(), ""})
		return
	}
	post, err := models.Post{ID: oid}.Get()
	if err != nil {
		c.JSON(http.StatusNotFound, response{false, err.Error(), "no post with that id"})
		return
	}
	var vote *models.Vote
	if err := c.Bind(vote); err != nil {
		c.JSON(http.StatusBadRequest, response{false, err.Error(), "invalid vote object"})
		return
	}
	if err := post.Vote(vote); err != nil {
		c.JSON(http.StatusInternalServerError, response{false, err.Error(), ""})
		return
	}
	c.JSON(http.StatusOK, response{true, post, "succesfully voted on post"})
}
