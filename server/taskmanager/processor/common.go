package processor

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"choccy/server/util"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func CheckAndRemoveUnValidDatabase(task *model.Task, databaseName string) {
	settingPath, err := util.GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	databasePath := filepath.Join(settingPath.CodeQLDatabase, databaseName)
	_, err = os.Stat(databasePath)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		panic(err.Error())
	}

	WriteTaskLog(task, "Check that the database is valid:"+databaseName)

	databaseYml := filepath.Join(databasePath, "codeql-database.yml")
	_, err = os.Stat(databaseYml)
	if os.IsNotExist(err) {
		WriteTaskLog(task, "There is no codeql-database in the database directory.yml, delete the invalid database:"+databasePath)
		os.RemoveAll(databasePath)
		return
	} else if err != nil {
		panic(err.Error())
	}

	fileBytes, err := os.ReadFile(databaseYml)
	if err != nil {
		panic(err.Error())
	}
	match := regexp.MustCompile(`finalised\s*:\s*(.+)`).FindSubmatch(fileBytes)
	if len(match) == 0 {
		WriteTaskLog(task, "codeql-database.The finalised field is not included in yml. Delete the invalid database:"+databasePath)
		os.RemoveAll(databasePath)
		return
	} else {
		if strings.TrimSpace(string(match[1])) == "false" {
			WriteTaskLog(task, "codeql-database.The finalised value in yml is false, and the invalid database is deleted:"+databasePath)
			os.RemoveAll(databasePath)
			return
		}
	}
}

func WriteTaskLog(task *model.Task, log string) {
	task.Logs += fmt.Sprintf("[%s] - %s\n", time.Now().Format("2006-01-02 15:04:05"), log)
	result := database.DB.Save(task)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}

func SetTaskStatus(task *model.Task, stats int) {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	task.Status = stats //0-in-queue, 1-in-progress, 2=completion, -1-error
	task.IsRead = false

	if stats == 2 {
		if setting.AutoReadCompletedTask {
			task.IsRead = true
		}
		if setting.AutoReadNoResultTask {
			if task.TotalResultsCount == 0 {
				task.IsRead = true
			}
		}
		if setting.AutoReadEmptyTask {
			if len(task.AnalyzedVersions) == 0 {
				task.IsRead = true
			}
		}
	}

	result = database.DB.Save(task)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}

func SetTaskStage(task *model.Task, stage int) {
	task.Stage = stage //0-Determine whether there is a new version, 1-Download the new version, 2-Compile the database, 3-Analysis
	task.IsRead = false
	result := database.DB.Save(task)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}

func SetTaskVersions(task *model.Task, versions []string) {
	task.Versions = versions
	result := database.DB.Save(task)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}

func AddTaskTotalResultsCount(task *model.Task, count int) {
	task.TotalResultsCount += count
	result := database.DB.Save(task)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}

func AddTaskAnalyzedVersion(task *model.Task, version string) {
	task.AnalyzedVersions = append(task.AnalyzedVersions, version)
	result := database.DB.Save(task)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}

func SetProjectLastAnalyzeReleaseTag(project *model.Project, tag string) {
	project.LastAnalyzeReleaseTag = tag
	project.LastAnalyzeTime = time.Now()
	result := database.DB.Save(project)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}

func SetProjectLastAnalyzeDatabaseCommit(project *model.Project, commit string) {
	project.LastAnalyzeDatabaseCommit = commit
	project.LastAnalyzeTime = time.Now()
	result := database.DB.Save(project)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}

func SetProjectLatestVersion(project *model.Project, version string, updateTime time.Time) {
	project.LatestVersion = version
	project.LatestVersionUpdateTime = updateTime
	project.LatestVersionCheckTime = time.Now()
	project.LatestVersionCheckMode = project.Mode
	result := database.DB.Save(project)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}

func CreateTaskResult(version string, commit string, fileName string, resultCount int, taskId uint) {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	taskResult := model.TaskResult{}
	taskResult.Version = version
	taskResult.Commit = commit
	taskResult.FileName = fileName
	taskResult.ResultCount = resultCount

	taskResult.IsRead = false
	taskResult.TaskId = taskId

	if setting.AutoReadNoResultResult {
		if taskResult.ResultCount == 0 {
			taskResult.IsRead = true
		}
	}

	result = database.DB.Save(&taskResult)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}

func Analyze(task *model.Task, databasePath, version string) (string, string) {
	WriteTaskLog(task, "Build a temporary query suite containing queries:"+strings.Join(task.ProjectSuite, " "))
	tmpSuitePath, err := util.GenerateSuite(task.ProjectSuite)
	if err != nil {
		panic("Failed to build temporary query suite:" + err.Error())
	}
	WriteTaskLog(task, "Temporary query suite build completed:"+tmpSuitePath)
	defer func() {
		WriteTaskLog(task, "Clean up the temporary query suite:"+tmpSuitePath)
		os.Remove(tmpSuitePath)
	}()
	WriteTaskLog(task, "Start analysis")

	resultFileName := fmt.Sprintf("%s__%s__%s__%s__%d.sarif", task.ProjectOwner, task.ProjectRepo, task.ProjectLanguage, version, time.Now().Unix())
	_, stderr, err, outputPath := util.DatabaseAnalyze(
		databasePath,
		tmpSuitePath,
		resultFileName,
	)
	if err != nil {
		panic("Analysis failed:" + stderr)
	}
	//writeTaskLog(task, "Scan log: "+stderr)
	WriteTaskLog(task, "Analysis completed, analysis results:"+outputPath)
	return resultFileName, outputPath
}
