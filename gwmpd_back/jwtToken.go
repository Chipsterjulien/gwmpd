package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func authenticator(userID string, password string, c *gin.Context) (interface{}, bool) {
	if userID == viper.GetString("ginserver.login") && password == viper.GetString("ginserver.password") {
		return userID, true
	}

	return userID, false
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
