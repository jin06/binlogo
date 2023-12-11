package user

import "github.com/gin-gonic/gin"

type AuthType string

const (
	Basic = "basic"
	Ldap  = "ldap"
)

type Handler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Info(c *gin.Context)
}
