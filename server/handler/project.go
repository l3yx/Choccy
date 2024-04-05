package handler

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"choccy/server/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func GetProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	sortBy := c.DefaultQuery("sortBy", "CreatedAt")
	switch sortBy {
	case "UpdatedAt":
		sortBy = "updated_at"
	case "CreatedAt":
		sortBy = "created_at"
	case "Url":
		sortBy = "url"
	case "Language":
		sortBy = "language"
	case "Mode":
		sortBy = "mode"
	case "Pause":
		sortBy = "pause"
	case "LatestVersionUpdateTime":
		sortBy = "latest_version_update_time"
	case "LastAnalyzeTime":
		sortBy = "last_analyze_time"
	default:
		sortBy = "unknown"
	}
	sortOrder := c.DefaultQuery("sortOrder", "descending")
	if sortOrder == "descending" {
		sortOrder = "desc"
	} else {
		sortOrder = "asc"
	}

	var projects []model.Project
	result := database.DB.
		Order(sortBy + " " + sortOrder).
		Scopes(database.Paginate(page, pageSize)).Find(&projects)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	var total int64
	database.DB.Model(&model.Project{}).Count(&total)

	//目的是获取更新检测时间间隔
	var setting model.Setting
	result = database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	//最大并发20，去获取最新版本并更新到数据库
	var wg sync.WaitGroup
	ch := make(chan int, 20)
	for index, _ := range projects {
		wg.Add(1)
		ch <- 0
		go func(i int) {
			defer func() { <-ch }()
			defer wg.Done()

			expire := projects[i].LatestVersionCheckTime.Add(time.Duration(float32(time.Minute) * setting.UpdateDetectionInterval))
			if expire.Before(time.Now()) || projects[i].LatestVersion == "[Error]" || projects[i].LatestVersionCheckMode != projects[i].Mode {
				if projects[i].Mode == 0 { //Release
					release, err := util.GetGithubReleaseLatest(projects[i].Owner, projects[i].Repo)
					if err != nil {
						log.Println("Error: " + err.Error())
						projects[i].LatestVersion = "[Error]"
						projects[i].LatestVersionErrorInfo = err.Error()
					} else {
						projects[i].LatestVersion = release.TagName
						projects[i].LatestVersionUpdateTime = release.CreatedAt
					}
				} else if projects[i].Mode == 1 { //原有数据库
					githubDatabase, err := util.GetGithubDatabase(projects[i].Owner, projects[i].Repo, projects[i].Language)
					if err != nil {
						log.Println("Error: " + err.Error())
						projects[i].LatestVersion = "[Error]"
						projects[i].LatestVersionErrorInfo = err.Error()
					} else {
						projects[i].LatestVersion = githubDatabase.CommitOid
						projects[i].LatestVersionUpdateTime = githubDatabase.CreatedAt
					}
				} else if projects[i].Mode == 3 { //默认分支
					branch, err := util.GetGithubDefaultBranch(projects[i].Owner, projects[i].Repo)
					if err != nil {
						log.Println("Error: " + err.Error())
						projects[i].LatestVersion = "[Error]"
						projects[i].LatestVersionErrorInfo = err.Error()
					} else {
						projects[i].LatestVersion = branch.Commit.Sha
						projects[i].LatestVersionUpdateTime = branch.Commit.Commit.Committer.Date
					}
				} else {
					projects[i].LatestVersion = "[Error]"
					projects[i].LatestVersionErrorInfo = "未知扫描模式"
				}
				projects[i].LatestVersionCheckTime = time.Now()
				projects[i].LatestVersionCheckMode = projects[i].Mode
				database.DB.Save(projects[i])
			}

		}(index)
	}
	wg.Wait()

	c.JSON(200, gin.H{
		"data":  projects,
		"total": total,
	})
}

func SaveProject(c *gin.Context) {
	var project model.Project
	err := c.ShouldBind(&project)
	if err != nil {
		panic(err.Error())
	}

	urlReg := regexp.MustCompile("^https://github\\.com/([^/]+)/([^/]+)$")
	matches := urlReg.FindStringSubmatch(project.Url)
	if len(matches) != 3 {
		panic(fmt.Sprintf("请填写正确项目地址: %s", urlReg))
	}
	owner := matches[1] // 提取 owner 字段
	repo := matches[2]  // 提取 repo 字段
	project.Owner = owner
	project.Repo = repo

	langReg := regexp.MustCompile("^[a-zA-Z]+$")
	if !langReg.MatchString(project.Language) {
		panic(fmt.Sprintf("请填写正确项目语言: %s", langReg))
	}

	if len(project.Suite) == 0 {
		panic(fmt.Sprintf("请选择查询套件"))
	}

	var result *gorm.DB
	if project.ID == 0 {
		result = database.DB.Save(&project)
	} else {
		result = database.DB.Model(&project).
			Select("Url", "Owner", "Repo", "Mode", "Language", "Command", "Suite", "Pause").
			Updates(project)
	}

	if result.Error != nil {
		panic(result.Error.Error())
	}

	c.JSON(200, gin.H{
		"data": project,
	})
}

func DeleteProject(c *gin.Context) {
	ID, err := strconv.Atoi(c.DefaultQuery("ID", ""))
	if err != nil {
		panic(err.Error())
	}
	result := database.DB.Delete(&model.Project{}, ID)
	if result.Error != nil {
		panic(result.Error.Error())
	}
	c.JSON(200, gin.H{
		"data": ID,
	})
}
