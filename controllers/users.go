package controllers

import (
	"dedpidgon/thoughtsea/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type users struct{}

func (x users) create(c *gin.Context) {

	var user *models.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, response{false, err.Error(), ""})
		return
	}

	if err := user.Save(); err != nil {
		c.JSON(http.StatusBadRequest, response{false, err.Error(), ""})
		return
	}

	c.JSON(http.StatusCreated, response{true, user, "Succesfully created user!"})
}

func (x users) read(c *gin.Context) {

	list, err := models.User{}.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{false, err.Error(), ""})
		return
	}
	c.JSON(http.StatusOK, response{true, list, "here are all the users"})
}

func (x users) readOne(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response{false, err.Error(), ""})
		return
	}
	user, err := models.User{ID: oid}.Get()
	if err != nil {
		c.JSON(http.StatusBadRequest, response{false, err.Error(), ""})
		return
	}
	c.JSON(http.StatusOK, response{true, user, "here is your user"})
}

// readAuthenticated gets the user from the token in the request header and returns it
func (x users) readAuthenticated(c *gin.Context) {

	headerValue, exists := c.Get(AUTHORIZATION)
	if !exists {
		// handle case where header does not exist
		return
	}
	oid, ok := headerValue.(primitive.ObjectID)
	if !ok {
		// handle case where header value is not a valid ObjectID
		return
	}

	user, err := models.User{ID: oid}.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{false, err.Error(), ""})
		return
	}
	c.JSON(http.StatusOK, response{true, user, "here is the currently logged in user"})
}

func (x users) update(c *gin.Context) {
	var user *models.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, response{false, err.Error(), ""})
		return
	}
	if !usersMatch(c, user.ID) {
		c.JSON(http.StatusUnauthorized, response{false, "you can not update another user", "how did you even get here?"})
		return
	}
	if err := user.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, response{false, err.Error(), ""})
		return
	}

	c.JSON(http.StatusOK, response{true, user, "your user has been updated, here is the new one"})
}

func (x users) destroy(c *gin.Context) {

	fmt.Println("in the user destroy method, going to try to get oid from the url")
	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		fmt.Println("failed to convert the id")
		c.JSON(http.StatusBadRequest, response{false, err.Error(), ""})
		return
	}

	fmt.Println("got the oid...")
	user := &models.User{ID: oid}
	fmt.Println("made the user...")
	if !usersMatch(c, user.ID) {
		c.JSON(http.StatusUnauthorized, response{false, "you can not delete another user", "shame on you..."})
		return
	}
	if err := user.Destroy(); err != nil {
		c.JSON(http.StatusInternalServerError, response{false, err.Error(), ""})
		return
	}
	c.JSON(http.StatusOK, response{true, user, "this user has been deleted"})
}

func (x users) login(c *gin.Context) {

	var user *models.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, response{false, err.Error(), "what even is that?"})
		return
	}

	if err := user.Login(); err != nil {
		c.JSON(http.StatusBadRequest, response{false, err.Error(), "can not log in with those creds"})
		return
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{false, err.Error(), "error generating the token"})
		return
	}

	c.JSON(http.StatusOK, response{true, token, "succesfully logged in, here is your token"})
}
