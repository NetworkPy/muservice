package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

// Create an account group
// Create a handler (which will later have injected services)
func NewHandler(c *Config) {
	h := &Handler{}
	g := c.Router.Group(os.Getenv("ACCOUNT_API_URL")) // Init group
	{
		g.GET("/me", h.Me)
		g.POST("/signup", h.Signup)
		g.POST("/signin", h.Signin)
		g.POST("/signout", h.Signout)
		g.POST("/tokens", h.Tokens)
		g.POST("/image", h.Image)
		g.DELETE("/image", h.DeleteImage)
		g.PUT("/details", h.Details)
	}
}

// Me handler calls services for getting
// a user's details
func (h *Handler) Me(cnx *gin.Context) {
	cnx.JSON(http.StatusOK, gin.H{
		"hello": "it's me",
	})
}

// Signup handler
func (h *Handler) Signup(cnx *gin.Context) {
	cnx.JSON(http.StatusOK, gin.H{
		"hello": "it's signup",
	})
}

// Signin handler
func (h *Handler) Signin(cnx *gin.Context) {
	cnx.JSON(http.StatusOK, gin.H{
		"hello": "it's signin",
	})
}

// Signout handler
func (h *Handler) Signout(cnx *gin.Context) {
	cnx.JSON(http.StatusOK, gin.H{
		"hello": "it's signout",
	})
}

// Tokens router
func (h *Handler) Tokens(cnx *gin.Context) {
	cnx.JSON(http.StatusOK, gin.H{
		"hello": "it's tokens",
	})
}

// Image router
func (h *Handler) Image(cnx *gin.Context) {
	cnx.JSON(http.StatusOK, gin.H{
		"hello": "it's signin",
	})
}

// DeleteImage handler
func (h *Handler) DeleteImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's deleteImage",
	})
}

// Details handler
func (h *Handler) Details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's details",
	})
}
