package controllers

import (
	"github.com/gin-gonic/gin"
)

const (
	USERS    string = "/users"
	POSTS    string = "/posts"
	COMMENTS string = "/comments"
)

// Core controller handles all client requests
type Core struct {
	*gin.Engine
}

// Init stores all the routes
func (c Core) Init() {
	// users
	c.POST(USERS, users{}.create)
	c.GET(USERS, users{}.read)
	c.GET(USERS+"/:id", users{}.readOne)
	c.GET(USERS+"/auth", AuthMiddleware(), users{}.readAuthenticated)
	c.PUT(USERS, AuthMiddleware(), users{}.update)
	c.DELETE(USERS+"/:id", AuthMiddleware(), users{}.destroy)
	c.POST(USERS+"/auth", users{}.login)

	// posts
	c.POST("/posts", AuthMiddleware(), posts{}.create)
	c.GET("/posts", posts{}.read)
	c.GET("/posts/:id", posts{}.readOne)
	c.PUT("/posts", AuthMiddleware(), posts{}.update)
	c.DELETE("/posts/:id", AuthMiddleware(), posts{}.destroy)
	c.POST("/posts/:id/vote", AuthMiddleware(), posts{}.vote)

	// comments
	c.POST("/comments", AuthMiddleware(), comments{}.create)
	c.GET("/comments", comments{}.read)
	c.GET("/comments/:id", comments{}.readOne)
	c.PUT("/comments", AuthMiddleware(), comments{}.update)
	c.DELETE("/comments/:id", AuthMiddleware(), comments{}.destroy)
	c.POST("/comments/:id/vote", AuthMiddleware(), comments{}.vote)

}
