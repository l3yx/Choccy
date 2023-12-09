package handler

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"choccy/server/taskmanager"
	"choccy/server/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
	"strconv"
	"strings"
)

func RunTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))

	var project model.Project
	result := database.DB.First(&project, id)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	success, err := taskmanager.AddTask(&project, true)
	if err != nil {
		panic(err.Error())
	}

	if success {
		c.JSON(200, gin.H{
			"data": gin.H{
				"ok": true,
			},
		})
	} else {
		c.JSON(200, gin.H{
			"data": gin.H{
				"ok": false,
			},
		})
	}
}

func GetTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	sortBy := c.DefaultQuery("sortBy", "UpdatedAt")
	switch sortBy {
	case "UpdatedAt":
		sortBy = "updated_at"
	case "ProjectID":
		sortBy = "project_id"
	case "Status":
		sortBy = "status"
	case "Stage":
		sortBy = "stage"
	case "CreatedAt":
		sortBy = "created_at"
	case "ProjectMode":
		sortBy = "project_mode"
	case "ProjectLanguage":
		sortBy = "project_language"
	case "ProjectName":
		sortBy = "project_name"
	case "TotalResultsCount":
		sortBy = "total_results_count"
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

	var tasks []model.Task
	result := database.DB.
		Where(filters).
		Order(sortBy + " " + sortOrder).
		Scopes(database.Paginate(page, pageSize)).
		Find(&tasks)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	var total int64
	result = database.DB.Model(&model.Task{}).Where(filters).Count(&total)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	c.JSON(200, gin.H{
		"data":  tasks,
		"total": total,
	})
}

func SetTaskIsRead(c *gin.Context) {
	var isRead model.IsRead
	err := c.ShouldBind(&isRead)
	if err != nil {
		panic(err.Error())
	}

	for _, id := range isRead.IdList {
		var task model.Task
		result := database.DB.First(&task, id)
		if result.Error != nil {
			panic(result.Error.Error())
		}
		task.IsRead = isRead.Read
		result = database.DB.Save(task)
		if result.Error != nil {
			panic(result.Error.Error())
		}
	}

	c.JSON(200, gin.H{
		"data": "ok",
	})
}

func GetTaskUnread(c *gin.Context) {
	var count int64
	result := database.DB.Model(&model.Task{}).Where("is_read = false").Count(&count)
	if result.Error != nil {
		panic(result.Error.Error())
	}
	c.JSON(200, gin.H{
		"count": count,
	})
}

func AddTask(c *gin.Context) {
	type Task struct {
		Database string   `json:"database"`
		Suites   []string `json:"suites"`
		Name     string   `json:"name"`
	}
	var task Task
	err := c.ShouldBind(&task)
	if err != nil {
		panic(err.Error())
	}

	if strings.TrimSpace(task.Database) == "" {
		panic(fmt.Sprintf("请选择数据库"))
	}
	if len(task.Suites) == 0 {
		panic(fmt.Sprintf("请选择查询套件"))
	}
	if strings.TrimSpace(task.Name) == "" {
		task.Name = task.Database
	}

	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	success, err := taskmanager.AddCustomTask(path.Join(settingPath.CodeQLDatabase, task.Database), task.Suites, task.Name)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"data": gin.H{
			"success": success,
		},
	})
}
