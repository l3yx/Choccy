package taskmanager

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"choccy/server/taskmanager/processor"
	"choccy/server/util"
	"fmt"
	"log"
	"os"
	"strings"
)

func Consumer() {
	for id := range CH {
		func() {
			var task model.Task
			result := database.DB.First(&task, id)
			if result.Error != nil {
				log.Println("Error: " + result.Error.Error())
				return
			}

			//异常处理
			defer func(task *model.Task) {
				if r := recover(); r != nil {
					processor.WriteTaskLog(task, fmt.Sprintf("%s", r))
					processor.SetTaskStatus(task, -1) //任务错误
				}
			}(&task)

			resultCount := 0

			processor.SetTaskStatus(&task, 1) //任务进行中

			var project model.Project
			result = database.DB.First(&project, task.ProjectID)
			if result.Error != nil {
				panic(fmt.Sprintf("获取项目 %d 失败", task.ProjectID))
			}

			var modelStr string
			if task.ProjectMode == 0 {
				modelStr = "Release"
			} else {
				modelStr = "原有数据库"
			}
			processor.WriteTaskLog(&task,
				fmt.Sprintf("开始任务，项目：%s/%s，语言：%s， 模式：%s，查询套件：%s",
					task.ProjectOwner, task.ProjectRepo,
					task.ProjectLanguage,
					modelStr,
					strings.Join(task.ProjectSuite, " "),
				),
			)

			if task.ProjectMode == 0 { //Release
				// 新版判断
				processor.SetTaskStage(&task, 0) //新版本判断
				tags, latestRelease := processor.CheckReleaseUpdates(&task, project.LastAnalyzeReleaseTag, &project)
				if len(tags) == 0 {
					if !task.Manual {
						processor.WriteTaskLog(&task, "当前没有新版本，结束任务")
						processor.SetTaskStatus(&task, 2) //任务完成
						return
					} else {
						processor.WriteTaskLog(&task, "当前没有新版本，但该任务手动触发，默认扫描当前最新版："+latestRelease.TagName)
						tags = []string{latestRelease.TagName}
					}
				} else {
					processor.WriteTaskLog(&task, "获取到新版本："+strings.Join(tags, "，"))
				}

				processor.SetTaskVersions(&task, tags)

				for _, tag := range tags {
					databaseName := fmt.Sprintf("%s__%s__%s__r__%s",
						task.ProjectOwner,
						task.ProjectRepo,
						task.ProjectLanguage,
						tag)
					githubTag, err := util.GetGithubTag(project.Owner, project.Repo, tag)
					if err != nil {
						panic("获取tag对应的commit失败：" + err.Error())
					}
					processor.CheckAndRemoveUnValidDatabase(&task, databaseName)
					databasePath := util.IsCodeQLDatabaseExists(databaseName)
					if databasePath == "" {
						//下载新版本
						processor.SetTaskStage(&task, 1) // 下载新版本
						tagSourcePath := processor.DownloadRelease(&task, tag)
						defer func() {
							processor.WriteTaskLog(&task, "清理代码："+tagSourcePath)
							os.RemoveAll(tagSourcePath)
						}()

						//编译数据库
						processor.SetTaskStage(&task, 2) // 编译数据库
						databasePath = processor.CreateDatabase(&task, tagSourcePath, databaseName)
					} else {
						processor.WriteTaskLog(&task, fmt.Sprintf("数据库 %s 有效，跳过源码下载和数据库构建", databaseName))
					}

					//扫描
					processor.SetTaskStage(&task, 3)
					resultFileName, resultFilePath := processor.Analyze(&task, databasePath, tag)
					codeQLSarif, err := util.ParseSarifFile(resultFilePath)
					if err != nil {
						panic("分析结果解析错误：" + err.Error())
					}
					resultCount += len(codeQLSarif.Results)
					processor.AddTaskTotalResultsCount(&task, len(codeQLSarif.Results))
					processor.WriteTaskLog(&task, fmt.Sprintf("扫描结果数量：%d", len(codeQLSarif.Results)))
					processor.AddTaskAnalyzedVersion(&task, tag)
					processor.SetProjectLastAnalyzeReleaseTag(&project, tag)
					processor.CreateTaskResult(tag, githubTag.Commit.Sha, resultFileName, len(codeQLSarif.Results), task.ID)
				}
			} else if task.ProjectMode == 1 { //原有数据库
				// 新版判断
				processor.SetTaskStage(&task, 0) //新版本判断
				databaseCommit, databaseUrl := processor.CheckDatabaseUpdates(&task, &project)
				if databaseCommit == project.LastAnalyzeDatabaseCommit {
					if !task.Manual {
						processor.WriteTaskLog(&task, "当前没有新版本，结束任务")
						processor.SetTaskStatus(&task, 2) //任务完成
						return
					} else {
						processor.WriteTaskLog(&task, "当前没有新版本，但该任务手动触发，默认扫描当前最新版："+databaseCommit)
					}
				} else {
					processor.WriteTaskLog(&task, "获取到新版本："+databaseCommit)
				}
				processor.SetTaskVersions(&task, []string{databaseCommit})
				databaseName := fmt.Sprintf("%s__%s__%s__d__%s",
					task.ProjectOwner,
					task.ProjectRepo,
					task.ProjectLanguage,
					databaseCommit[:7])
				processor.CheckAndRemoveUnValidDatabase(&task, databaseName)
				databasePath := util.IsCodeQLDatabaseExists(databaseName)
				if databasePath == "" {
					// 下载新版本
					processor.SetTaskStage(&task, 1) // 下载新版本
					databasePath = processor.DownloadDatabase(&task, databaseUrl, databaseCommit, databaseName)
				} else {
					processor.WriteTaskLog(&task, fmt.Sprintf("数据库 %s 有效，跳过数据库下载", databaseName))
				}

				// 扫描
				processor.SetTaskStage(&task, 3) // 扫描
				resultFileName, resultFilePath := processor.Analyze(&task, databasePath, databaseCommit[:7])
				codeQLSarif, err := util.ParseSarifFile(resultFilePath)
				if err != nil {
					panic("分析结果解析错误：" + err.Error())
				}
				resultCount += len(codeQLSarif.Results)
				processor.AddTaskTotalResultsCount(&task, len(codeQLSarif.Results))
				processor.WriteTaskLog(&task, fmt.Sprintf("扫描结果数量：%d", len(codeQLSarif.Results)))
				processor.AddTaskAnalyzedVersion(&task, databaseCommit)
				processor.SetProjectLastAnalyzeDatabaseCommit(&project, databaseCommit)
				processor.CreateTaskResult(databaseCommit, databaseCommit, resultFileName, len(codeQLSarif.Results), task.ID)
			}

			processor.SetTaskStatus(&task, 2) //任务完成
		}()
	}
}
