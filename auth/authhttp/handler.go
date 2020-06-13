package authhttp

import (
	"fmt"
	"net/http"

	"github.com/erickoliv/finances-api/auth"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	auth   auth.Authenticator
	signer auth.SessionSigner
}

type credential struct {
	Username string `json:"username" binding:"required" `
	Password string `json:"password" binding:"required" `
}

// NewHTTPHandler receives a authenticator and a signer Service, returning a internal a HTTP account handler
func NewHTTPHandler(authenticator auth.Authenticator, signer auth.SessionSigner) *HTTPHandler {
	return &HTTPHandler{
		auth:   authenticator,
		signer: signer,
	}
}

func (handler *HTTPHandler) Router(group *gin.RouterGroup) {
	group.POST("/login", handler.login)
	group.POST("/register", handler.register)
}

func (handler *HTTPHandler) login(c *gin.Context) {
	credentials := credential{}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid payload",
		})
		return
	}

	user, err := handler.auth.Login(c, credentials.Username, credentials.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	cookie, err := handler.signer.SignUser(user.UUID.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.SetCookie(auth.AuthCookie, cookie, 0, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": cookie,
	})
}

func (handler *HTTPHandler) register(c *gin.Context) {
	user := &auth.User{}
	if err := c.Bind(user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := handler.auth.Register(c, user); err != nil {
		fmt.Printf("error to register user %v - %v \n", user, err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "registration error",
		})
		return
	}

	c.Status(http.StatusCreated)
}
