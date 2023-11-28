package handler

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"choccy/server/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strconv"
)

func GetResult(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	sortBy := c.DefaultQuery("sortBy", "CreatedAt")
	switch sortBy {
	case "CreatedAt":
		sortBy = "created_at"
	case "ResultCount":
		sortBy = "result_count"
	case "FileName":
		sortBy = "file_name"
	default:
		sortBy = "unknown"
	}
	sortOrder := c.DefaultQuery("sortOrder", "descending")
	if sortOrder == "descending" {
		sortOrder = "desc"
	} else {
		sortOrder = "asc"
	}
	filtersStr := c.DefaultQuery("filters", "{}")
	var filters map[string]interface{}
	err := json.Unmarshal([]byte(filtersStr), &filters)
	if err != nil {
		panic(err.Error())
	}

	var taskResults []model.TaskResult
	result := database.DB.
		Preload("Task").
		Where(filters).
		Order(sortBy + " " + sortOrder).
		Scopes(database.Paginate(page, pageSize)).
		Find(&taskResults)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	var total int64
	result = database.DB.Model(&model.TaskResult{}).Where(filters).Count(&total)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	for index, _ := range taskResults {
		resultFilePath := filepath.Join(settingPath.CodeQlResult, taskResults[index].FileName)
		taskResults[index].Task.Logs = ""
		taskResults[index].FilePath = resultFilePath
	}

	c.JSON(200, gin.H{
		"data":  taskResults,
		"total": total,
	})
}

func GetResultSarif(c *gin.Context) {
	fileName := c.DefaultQuery("fileName", "unknow")
	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}
	resultFilePath := filepath.Join(settingPath.CodeQlResult, fileName)
	codeQLSarif, err := util.ParseSarifFile(resultFilePath)
	if err != nil {
		codeQLSarif = &util.CodeQLSarif{}
		codeQLSarif.NotificationsId = []string{}
		codeQLSarif.Rules = []string{}
		codeQLSarif.Packs = []string{}
		codeQLSarif.Results = make([]util.CodeQLResult, 0)
	}
	c.JSON(200, gin.H{
		"data": codeQLSarif,
	})
}

func DeleteResult(c *gin.Context) {
	ID, err := strconv.Atoi(c.DefaultQuery("ID", ""))
	if err != nil {
		panic(err.Error())
	}
	var taskResult model.TaskResult

	result := database.DB.First(&taskResult, ID)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	resultFilePath, err := filepath.Abs(filepath.Join(settingPath.CodeQlResult, taskResult.FileName))
	if err != nil {
		panic(err.Error())
	}
	if filepath.Dir(resultFilePath) != settingPath.CodeQlResult {
		panic("Wrong file name")
	}

	result = database.DB.Delete(&taskResult)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	_ = os.Remove(resultFilePath)

	c.JSON(200, gin.H{
		"data": resultFilePath,
	})
}

func SetResultIsRead(c *gin.Context) {
	var isRead model.IsRead
	err := c.ShouldBind(&isRead)
	if err != nil {
		panic(err.Error())
	}

	for _, id := range isRead.IdList {
		var taskResult model.TaskResult
		result := database.DB.First(&taskResult, id)
		if result.Error != nil {
			panic(result.Error.Error())
		}
		taskResult.IsRead = isRead.Read
		result = database.DB.Save(taskResult)
		if result.Error != nil {
			panic(result.Error.Error())
		}
	}

	c.JSON(200, gin.H{
		"data": "ok",
	})
}

func GetResultUnread(c *gin.Context) {
	var count int64
	result := database.DB.Model(&model.TaskResult{}).Where("is_read = false").Count(&count)
	if result.Error != nil {
		panic(result.Error.Error())
	}
	c.JSON(200, gin.H{
		"count": count,
	})
}
