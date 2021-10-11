package handler

import (
	"net/http"
	"os"

	"github.com/NetworkPy/muserv/muservice/account/models"
	"github.com/gin-gonic/gin"
)

// Handler struct holds required services for handler to function
type Handler struct {
	UserService  models.UserService
	TokenService models.TokenService
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	Router       *gin.Engine
	UserService  models.UserService
	TokenService models.TokenService
}

// Create an account group
// Create a handler (which will later have injected services)
func NewHandler(c *Config) {
	h := &Handler{
		UserService:  c.UserService,
		TokenService: c.TokenService,
	}

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
