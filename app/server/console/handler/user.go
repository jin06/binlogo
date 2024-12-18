package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jin06/binlogo/v2/app/server/console/basic"
	"github.com/jin06/binlogo/v2/app/server/console/service"
	"github.com/spf13/viper"
)

type token struct {
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req loginReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(200, basic.Fail(err))
		return
	}
	pass := service.DefaultAuth().Verify(req.Username, req.Password)
	if !pass {
		c.JSON(200, basic.FailCode(basic.CodeBadUsernameOrPassword))
		return
	}
	t := service.DefaultStore().Set()
	c.JSON(200, basic.Success(map[string]string{"token": t}))
}

func Logout(c *gin.Context) {
	t := c.GetHeader("x-token")
	service.DefaultStore().Remove(t)
	c.JSON(200, basic.Success(nil))
}

func UserInfo(c *gin.Context) {
	t := c.Query("token")
	if service.DefaultStore().Get(t) {
		c.JSON(200, basic.Success(token{
			Roles:        []string{"admin"},
			Introduction: "I am a super administrator",
			Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			Name:         "Super Admin",
		}))
		return
	}
	c.JSON(200, basic.FailCode(basic.CodeTokenExpired))
}

func AuthType(c *gin.Context) {
	typ := viper.GetString("auth.authorizer.type")
	c.JSON(200, basic.Success(map[string]any{
		"type": typ,
	}))
}
