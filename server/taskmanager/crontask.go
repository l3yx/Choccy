package taskmanager

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"github.com/robfig/cron/v3"
	"log"
)

func runProjects() {
	var projects []model.Project
	result := database.DB.Where("pause = ?", false).Find(&projects)
	if result.Error != nil {
		log.Println("Error: " + result.Error.Error())
	}

	for _, project := range projects {
		_, err := AddTask(&project, false)
		if err != nil {
			log.Println("Error: " + err.Error())
		}
	}
}

var TaskCron *cron.Cron

func SetCronTask() error {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		return result.Error
	}

	_, err := cron.ParseStandard(setting.CronTaskSpec)
	if err != nil {
		return err
	}

	if TaskCron != nil {
		TaskCron.Stop()
	}
	TaskCron = cron.New()
	_, err = TaskCron.AddFunc(setting.CronTaskSpec, runProjects)
	TaskCron.Start()
	if err != nil {
		return err
	}

	return nil
}
