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

			//Exception handling
			defer func(task *model.Task) {
				if r := recover(); r != nil {
					processor.WriteTaskLog(task, fmt.Sprintf("%s", r))
					processor.SetTaskStatus(task, -1) //task error
				}
			}(&task)

			resultCount := 0

			processor.SetTaskStatus(&task, 1) //Mission in progress

			var project model.Project
			result = database.DB.First(&project, task.ProjectID)
			if result.Error != nil {
				panic(fmt.Sprintf("Get item %d failed", task.ProjectID))
			}

			var modelStr string
			if task.ProjectMode == 0 {
				modelStr = "Release"
			} else {
				modelStr = "original database"
			}
			processor.WriteTaskLog(&task,
				fmt.Sprintf("Start Task, Item:%s/%s, Language:%s, Mode:%s, Query Suite:%s",
					task.ProjectOwner, task.ProjectRepo,
					task.ProjectLanguage,
					modelStr,
					strings.Join(task.ProjectSuite, " "),
				),
			)

			if task.ProjectMode == 0 { //Release
				// new edition judgment
				processor.SetTaskStage(&task, 0) // New version judgment
				tags, latestRelease := processor.CheckReleaseUpdates(&task, project.LastAnalyzeReleaseTag, &project)
				if len(tags) == 0 {
					if !task.Manual {
						processor.WriteTaskLog(&task, "There is no new version, end the task")
						processor.SetTaskStatus(&task, 2) //Task completed
						return
					} else {
						processor.WriteTaskLog(&task, "There is currently no new version, but the task is triggered manually, and the current latest version is scanned by default:"+latestRelease.TagName)
						tags = []string{latestRelease.TagName}
					}
				} else {
					processor.WriteTaskLog(&task, "Get a new version:"+strings.Join(tags, ", "))
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
						panic("Failed to get commit corresponding to tag:" + err.Error())
					}
					processor.CheckAndRemoveUnValidDatabase(&task, databaseName)
					databasePath := util.IsCodeQLDatabaseExists(databaseName)
					if databasePath == "" {
						//Download new version
						processor.SetTaskStage(&task, 1) // Download new version
						tagSourcePath := processor.DownloadRelease(&task, tag)
						defer func() {
							processor.WriteTaskLog(&task, "Cleanup code:"+tagSourcePath)
							os.RemoveAll(tagSourcePath)
						}()

						//Compile database
						processor.SetTaskStage(&task, 2) // Compile database
						databasePath = processor.CreateDatabase(&task, tagSourcePath, databaseName)
					} else {
						processor.WriteTaskLog(&task, fmt.Sprintf("Database %s valid, skip source code download and database build", databaseName))
					}

					//Scanning
					processor.SetTaskStage(&task, 3)
					resultFileName, resultFilePath := processor.Analyze(&task, databasePath, tag)
					codeQLSarif, err := util.ParseSarifFile(resultFilePath)
					if err != nil {
						panic("Analysis result parsing error:" + err.Error())
					}
					resultCount += len(codeQLSarif.Results)
					processor.AddTaskTotalResultsCount(&task, len(codeQLSarif.Results))
					processor.WriteTaskLog(&task, fmt.Sprintf("Number of scan results: %d", len(codeQLSarif.Results)))
					processor.AddTaskAnalyzedVersion(&task, tag)
					processor.SetProjectLastAnalyzeReleaseTag(&project, tag)
					processor.CreateTaskResult(tag, githubTag.Commit.Sha, resultFileName, len(codeQLSarif.Results), task.ID)
				}
			} else if task.ProjectMode == 1 { //original database
				// new edition judgment
				processor.SetTaskStage(&task, 0) //new edition judgment
				databaseCommit, databaseUrl := processor.CheckDatabaseUpdates(&task, &project)
				if databaseCommit == project.LastAnalyzeDatabaseCommit {
					if !task.Manual {
						processor.WriteTaskLog(&task, "There is no new version, end the task")
						processor.SetTaskStatus(&task, 2) //Mission accomplished
						return
					} else {
						processor.WriteTaskLog(&task, "There is currently no new version, but the task is triggered manually, defaulting to the current latest version of Scanning:"+databaseCommit)
					}
				} else {
					processor.WriteTaskLog(&task, "Get a new version:"+databaseCommit)
				}
				processor.SetTaskVersions(&task, []string{databaseCommit})
				databaseCommitAbbr := databaseCommit
				if len(databaseCommit) > 7 {
					databaseCommitAbbr = databaseCommit[:7]
				}
				databaseName := fmt.Sprintf("%s__%s__%s__d__%s",
					task.ProjectOwner,
					task.ProjectRepo,
					task.ProjectLanguage,
					databaseCommitAbbr)
				processor.CheckAndRemoveUnValidDatabase(&task, databaseName)
				databasePath := util.IsCodeQLDatabaseExists(databaseName)
				if databasePath == "" {
					// Download new version
					processor.SetTaskStage(&task, 1) // Download new version
					databasePath = processor.DownloadDatabase(&task, databaseUrl, databaseCommit, databaseName)
				} else {
					processor.WriteTaskLog(&task, fmt.Sprintf("database %s valid, skip database download", databaseName))
				}

				// Scanning
				processor.SetTaskStage(&task, 3) // Scanning
				resultFileName, resultFilePath := processor.Analyze(&task, databasePath, databaseCommitAbbr)
				codeQLSarif, err := util.ParseSarifFile(resultFilePath)
				if err != nil {
					panic("Analysis result parsing error:" + err.Error())
				}
				resultCount += len(codeQLSarif.Results)
				processor.AddTaskTotalResultsCount(&task, len(codeQLSarif.Results))
				processor.WriteTaskLog(&task, fmt.Sprintf("Number of Scanning Results: %d", len(codeQLSarif.Results)))
				processor.AddTaskAnalyzedVersion(&task, databaseCommit)
				processor.SetProjectLastAnalyzeDatabaseCommit(&project, databaseCommit)
				processor.CreateTaskResult(databaseCommit, databaseCommit, resultFileName, len(codeQLSarif.Results), task.ID)
			}

			processor.SetTaskStatus(&task, 2) //Mission accomplished
		}()
	}
}
