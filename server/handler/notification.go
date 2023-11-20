package handler

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	var notifications []string

	var totalTasks int64
	var completedTasks int64
	var failedTasks int64
	var totalResults int64
	database.DB.Model(&model.Task{}).Count(&totalTasks)
	if totalTasks > 0 {
		database.DB.Model(&model.Task{}).Where("status = 2").Count(&completedTasks)
		database.DB.Model(&model.Task{}).Where("status = -1").Count(&failedTasks)
		database.DB.Model(&model.Task{}).Select("sum(total_results_count)").Scan(&totalResults)
	}

	var status model.Status
	result := database.DB.Take(&status)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	if totalTasks > status.TotalTasks {
		notifications = append(notifications, fmt.Sprintf("新增%d个任务", totalTasks-status.TotalTasks))
	}

	if completedTasks > status.CompletedTasks {
		notifications = append(notifications, fmt.Sprintf("%d个任务已执行完成", completedTasks-status.CompletedTasks))
	}

	if failedTasks > status.FailedTasks {
		notifications = append(notifications, fmt.Sprintf("%d个任务执行失败", failedTasks-status.FailedTasks))
	}

	if totalResults > status.TotalResults {
		notifications = append(notifications, fmt.Sprintf("扫到%d条数据", totalResults-status.TotalResults))
	}

	status.TotalTasks = totalTasks
	status.FailedTasks = failedTasks
	status.CompletedTasks = completedTasks
	status.TotalResults = totalResults
	database.DB.Save(status)

	c.JSON(200, gin.H{
		"data": gin.H{
			"notifications": notifications,
		},
	})
}
