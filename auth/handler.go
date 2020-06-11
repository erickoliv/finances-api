package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/erickoliv/finances-api/domain"
	"github.com/gin-gonic/gin"
)

type SessionSigner interface {
	SignUser(identifier string) (string, error)
	Validate(token string) (string, error)
}

type Authenticator interface {
	Login(ctx context.Context, username string, password string) (*domain.User, error)
	Register(ctx context.Context, user *domain.User) error
}

type HTTPHandler struct {
	auth   Authenticator
	signer SessionSigner
}

type credential struct {
	Username string `json:"username" binding:"required" `
	Password string `json:"password" binding:"required" `
}

// NewHTTPHandler receives a Account Service, returning a internal a HTTP account handler
func NewHTTPHandler(authenticator Authenticator, signer SessionSigner) *HTTPHandler {
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

	c.SetCookie(domain.AuthCookie, cookie, 0, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": cookie,
	})
}

func (handler *HTTPHandler) register(c *gin.Context) {
	user := &domain.User{}
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
