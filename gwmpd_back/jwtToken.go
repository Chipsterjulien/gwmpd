package main

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// func authenticator(userID string, password string, c *gin.Context) (interface{}, bool) {
func authenticator(c *gin.Context) (interface{}, error) {
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	if loginVals.Username == viper.GetString("ginserver.login") && loginVals.Password == viper.GetString("ginserver.password") {
		return loginVals.Username, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)

	return claims["id"].(string)
}

func payloadFunc(data interface{}) jwt.MapClaims {
	return jwt.MapClaims{identityKey: data}
}

func refreshToken(userID interface{}, c *gin.Context) bool {
	if userID == viper.GetString("ginserver.login") {
		return true
	}

	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
