package handler

import (
	"choccy/server/util"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"path/filepath"
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

func UploadDatabases(c *gin.Context) {
	file, _ := c.FormFile("file")

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	tmpDir, err := util.GetTmpDir()
	if err != nil {
		panic(err.Error())
	}

	dst := path.Join(tmpDir, file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		panic(err.Error())
	}
	defer os.Remove(dst)

	level, err := util.CheckDatabaseZip(dst)
	if err != nil {
		panic(err.Error())
	}

	if level == -1 || level > 1 {
		panic("无法识别数据库")
	}

	databasePath := path.Join(settingPath.CodeQLDatabase, file.Filename[:len(file.Filename)-len(filepath.Ext(file.Filename))])
	err = util.Unzip(dst, databasePath, level)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"data": databasePath,
	})
}

func DeleteDatabases(c *gin.Context) {
	databaseName := c.DefaultQuery("name", "")

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	databasePath, err := filepath.Abs(filepath.Join(settingPath.CodeQLDatabase, databaseName))
	if err != nil {
		panic(err.Error())
	}
	if filepath.Dir(databasePath) != settingPath.CodeQLDatabase {
		panic("文件名错误")
	}

	err = os.RemoveAll(databasePath)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"data": databasePath,
	})
}
