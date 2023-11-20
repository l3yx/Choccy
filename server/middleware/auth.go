package middleware

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}
	token := c.GetHeader("X-Token")
	if token == setting.SystemToken {
		c.Next()
	} else {
		c.JSON(200, gin.H{"err": "Unauthorized"})
		c.Abort()
	}
}
