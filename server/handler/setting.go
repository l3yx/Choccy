package handler

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"choccy/server/taskmanager"
	"choccy/server/util"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"os"
	"strings"
)

func GetSetting(c *gin.Context) {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}
	if taskmanager.TaskCron != nil {
		for _, entry := range taskmanager.TaskCron.Entries() {
			setting.CronTaskNextTime = entry.Next
			break
		}
	}
	c.JSON(200, gin.H{
		"data": setting,
	})
}

func SaveSetting(c *gin.Context) {
	var setting model.Setting
	err := c.ShouldBind(&setting)
	if err != nil {
		panic(err.Error())
	}

	if strings.TrimSpace(setting.SystemToken) == "" {
		panic("系统Token不能为空")
	}

	setting.CronTaskSpec = strings.TrimSpace(setting.CronTaskSpec)
	_, err = cron.ParseStandard(setting.CronTaskSpec)
	if err != nil {
		panic(err.Error())
	}

	setting.ID = 1
	result := database.DB.Save(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	err = taskmanager.SetCronTask()
	if err != nil {
		panic(err.Error())
	}
	util.SetEnv(setting.EnvStr)

	c.JSON(200, gin.H{
		"data": setting,
	})
}

func TestSetting(c *gin.Context) {
	var jsonData map[string]string
	err := c.ShouldBind(&jsonData)
	if err != nil {
		panic(err.Error())
	}

	key := jsonData["key"]
	value := jsonData["value"]
	var result interface{}

	if key == "CodeQLCli" {
		codeqlPath, err := util.GetCodeQL(value)
		if err != nil {
			panic(err.Error())
		}

		if value != "" {
			result = util.GetCodeQLVersionByPath(codeqlPath)
		} else {
			panic("CodeQL Cli cannot be an empty string")
		}

	} else if key == "EnvStr" {
		util.SetEnv(value)
		tmpResult := make(map[string]string)

		envs := os.Environ()
		for _, env := range envs {
			keyValue := strings.SplitN(env, "=", 2)
			if len(keyValue) == 2 {
				tmpResult[keyValue[0]] = keyValue[1]
			}
		}
		result = tmpResult
	}

	c.JSON(200, gin.H{
		"data": result,
	})
}
