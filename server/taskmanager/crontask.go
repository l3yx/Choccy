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

var c *cron.Cron
var Schedule cron.Schedule

func SetCronTask() error {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		return result.Error
	}

	schedule, err := cron.ParseStandard(setting.CronTaskSpec)
	if err != nil {
		return err
	}

	if c != nil {
		c.Stop()
	}
	c = cron.New()
	_, err = c.AddFunc(setting.CronTaskSpec, runProjects)
	c.Start()
	Schedule = schedule
	if err != nil {
		return err
	}

	return nil
}
