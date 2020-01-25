package auth

import "github.com/gin-gonic/gin"

// Routes related to auth package
func Routes(r *gin.RouterGroup) {
	r.POST("/register", Register)
	r.POST("/login", Login)
}
