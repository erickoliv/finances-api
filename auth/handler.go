package auth

import (
	"net/http"

	"github.com/erickoliv/finances-api/domain"
	"github.com/erickoliv/finances-api/service"
	"github.com/gin-gonic/gin"
)

// HTTPHandler is a interface to expose account manipulation using HTTP
type HTTPHandler interface {
	Router(*gin.RouterGroup)
}

type httpHandler struct {
	auth service.Authenticator
	sign service.Signer
}

type credential struct {
	Username string `json:"username" binding:"required" `
	Password string `json:"password" binding:"required" `
}

// NewHTTPHandler receives a Account Service, returning a internal a HTTP account handler
func NewHTTPHandler(authenticator service.Authenticator, signer service.Signer) HTTPHandler {
	return &httpHandler{
		auth: authenticator,
		sign: signer,
	}
}

func (handler *httpHandler) Router(group *gin.RouterGroup) {
	group.POST("/login", handler.login)
	group.POST("/register", handler.register)
}

func (handler *httpHandler) login(c *gin.Context) {

	credentials := credential{}
	if err := c.Bind(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid payload",
		})
		return
	}

	user, err := handler.auth.Login(c, credentials.Username, credentials.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	cookie, err := handler.sign.SignUser(c, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	c.SetCookie(domain.AuthCookie, cookie, 0, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": cookie,
	})
}

func (handler *httpHandler) register(c *gin.Context) {
	user := &domain.User{}
	if err := c.Bind(user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := handler.auth.Register(c, user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "registration error",
		})
		return
	}

	c.Status(http.StatusCreated)
}
