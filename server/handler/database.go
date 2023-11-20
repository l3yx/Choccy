package handler

import (
	"choccy/server/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetDatabases(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	sortBy := c.DefaultQuery("sortBy", "ModTime")
	sortOrder := c.DefaultQuery("sortOrder", "descending")

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	data, total, err := util.ListFiles(true, false, []string{"*"}, settingPath.CodeQLDatabase, sortBy, sortOrder, pageSize, page)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"data":  data,
		"total": total,
	})
}
