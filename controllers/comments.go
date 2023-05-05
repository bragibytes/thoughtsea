package controllers

import (
	"dedpidgon/thoughtsea/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type comments struct {
}

func (x comments) create(c *gin.Context) {
	var a *models.Comment
	if err := c.Bind(&a); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	oid, err := primitive.ObjectIDFromHex(c.GetHeader("Authorization"))
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
	}
	a.Author = oid
	if err := a.Save(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, a)
}
func (x comments) read(c *gin.Context) {
	list, err := models.Comment{}.GetAll()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}
func (x comments) readOne(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	comment, err := models.Comment{ID: oid}.Get()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, comment)
}
func (x comments) update(c *gin.Context) {
	var a *models.Comment
	if err := c.Bind(&a); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if !usersMatch(c, a.Author) {
		c.String(http.StatusUnauthorized, "i cant let you do that")
		return
	}
	if err := a.Update(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, a)
}
func (x comments) destroy(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	a := &models.Comment{ID: oid}
	if !usersMatch(c, a.Author) {
		c.String(http.StatusUnauthorized, "i cant let you do that")
		return
	}
	if err := a.Destroy(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, a)
}

func (x comments) vote(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response{false, err.Error(), ""})
		return
	}
	comment, err := models.Post{ID: oid}.Get()
	if err != nil {
		c.JSON(http.StatusNotFound, response{false, err.Error(), "no comment with that id"})
		return
	}
	var vote *models.Vote
	if err := c.Bind(vote); err != nil {
		c.JSON(http.StatusBadRequest, response{false, err.Error(), "invalid vote object"})
		return
	}
	if err := comment.Vote(vote); err != nil {
		c.JSON(http.StatusInternalServerError, response{false, err.Error(), ""})
		return
	}
	c.JSON(http.StatusOK, response{true, comment, "succesfully voted on post"})
}
