package handler

import (
	"choccy/server/util"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ResolveSuiteQueries(c *gin.Context) {
	path := c.DefaultQuery("path", "")
	if strings.TrimSpace(path) != "" {
		queries := util.ResolveSuiteQueries(path)
		var result []map[string]interface{}
		for _, query := range queries {
			metadata, err := util.ResolveQueryMetadata(query)
			if err != nil {
				panic(err.Error())
			}
			result = append(result, metadata)
		}
		c.JSON(200, gin.H{
			"data": result,
		})
	} else {
		panic("path cannot be empty")
	}
}

func SaveSuiteContent(c *gin.Context) {
	var jsonData map[string]string
	err := c.ShouldBind(&jsonData)
	if err != nil {
		panic(err.Error())
	}

	name := jsonData["name"]
	content := jsonData["content"]

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	filePath := filepath.Join(settingPath.CodeQLSuite, name)

	err = os.WriteFile(filePath, []byte(content), os.ModePerm)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"data": filePath,
	})
}

func GetSuiteContent(c *gin.Context) {
	name := c.DefaultQuery("name", "")

	if filepath.Ext(name) != ".qls" {
		panic("Must be a qls file")
	}

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	filePath := filepath.Join(settingPath.CodeQLSuite, name)

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"data": string(bytes),
	})
}

func GetSuites(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	sortBy := c.DefaultQuery("sortBy", "ModTime")
	sortOrder := c.DefaultQuery("sortOrder", "descending")

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	data, total, err := util.ListFiles(false, true, []string{".qls"}, settingPath.CodeQLSuite, sortBy, sortOrder, pageSize, page)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"data":  data,
		"total": total,
	})
}

func DeleteSuite(c *gin.Context) {
	suiteName := c.DefaultQuery("name", "")

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	suiteFilePath, err := filepath.Abs(filepath.Join(settingPath.CodeQLSuite, suiteName))
	if err != nil {
		panic(err.Error())
	}
	if filepath.Dir(suiteFilePath) != settingPath.CodeQLSuite {
		panic("File name error")
	}

	err = os.Remove(suiteFilePath)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"data": suiteFilePath,
	})
}

func CreateSuite(c *gin.Context) {
	suiteName := c.DefaultQuery("name", "")
	if !strings.HasSuffix(suiteName, ".qls") {
		suiteName = suiteName + ".qls"
	}

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	suiteFilePath, err := filepath.Abs(filepath.Join(settingPath.CodeQLSuite, suiteName))
	if err != nil {
		panic(err.Error())
	}
	if filepath.Dir(suiteFilePath) != settingPath.CodeQLSuite {
		panic("File name error")
	}

	_, err = os.Stat(suiteFilePath)
	if err == nil {
		panic("File already exists")
	}

	file, err := os.Create(suiteFilePath)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	_, err = file.WriteString("- description: new suite\n")
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"data": suiteFilePath,
	})
}

func RenameSuite(c *gin.Context) {
	oldSuiteName := c.DefaultQuery("oldName", "")
	newSuiteName := c.DefaultQuery("newName", "")
	if !strings.HasSuffix(newSuiteName, ".qls") {
		newSuiteName = newSuiteName + ".qls"
	}

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	oldSuiteFilePath, err := filepath.Abs(filepath.Join(settingPath.CodeQLSuite, oldSuiteName))
	if err != nil {
		panic(err.Error())
	}
	if filepath.Dir(oldSuiteFilePath) != settingPath.CodeQLSuite {
		panic("File name error")
	}

	newSuiteFilePath, err := filepath.Abs(filepath.Join(settingPath.CodeQLSuite, newSuiteName))
	if err != nil {
		panic(err.Error())
	}
	if filepath.Dir(newSuiteFilePath) != settingPath.CodeQLSuite {
		panic("File name error")
	}

	_, err = os.Stat(newSuiteFilePath)
	if err == nil {
		panic("File already exists")
	}

	err = os.Rename(oldSuiteFilePath, newSuiteFilePath)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"data": newSuiteFilePath,
	})
}
