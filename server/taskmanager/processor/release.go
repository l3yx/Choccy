package processor

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"choccy/server/util"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// CheckReleaseUpdates 返回新版本列表，理论上数组顺序是从旧到新
func CheckReleaseUpdates(task *model.Task, lastAnalyzeReleaseTag string, project *model.Project) ([]string, *util.GithubRelease) {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	releases, err := util.GetGithubReleases(task.ProjectOwner, task.ProjectRepo)
	if err != nil {
		panic("Failed to get project release:" + err.Error())
	}
	if len(releases) == 0 {
		panic("Project does not exist Release")
	}

	SetProjectLatestVersion(project, releases[0].TagName, releases[0].CreatedAt)

	var tags []string
	if strings.TrimSpace(lastAnalyzeReleaseTag) == "" {
		releaseCount := setting.FirstReleaseCount
		WriteTaskLog(task, fmt.Sprintf("The project scans the release for the first time, scans the latest %d versions, and subsequent scans will scan the increments.", releaseCount))
		for index, release := range releases {
			if index >= releaseCount {
				break
			}
			tags = append(tags, release.TagName)
		}
	} else {
		for _, release := range releases {
			if lastAnalyzeReleaseTag != release.TagName {
				tags = append(tags, release.TagName)
			} else {
				break
			}
		}
	}
	util.Reverse(&tags)

	return tags, &releases[0]
}

func DownloadRelease(task *model.Task, tag string) string {
	WriteTaskLog(task, "Download Version:"+tag)
	tagSourcePath, err := util.DownloadGithubTag(task.ProjectOwner, task.ProjectRepo, tag)
	if err != nil {
		panic("Download failed:" + err.Error())
	}
	WriteTaskLog(task, "Download successfully, path:"+tagSourcePath)
	return tagSourcePath
}

func CreateDatabase(task *model.Task, source string, databaseName string) string {
	WriteTaskLog(task, "Start building the database")
	stdout, stderr, err, databasePath := util.DatabaseCreate(
		source,
		task.ProjectLanguage,
		task.ProjectCommand,
		databaseName,
	)
	if err != nil {
		if !strings.Contains(stderr, "exists and is not an empty directory") {
			WriteTaskLog(task, "Clean up the database that failed to build:"+databasePath)
			os.RemoveAll(databasePath)
		}
		outError := ""
		stdoutLines := strings.Split(stdout, "\n")
		for _, stdoutLine := range stdoutLines {
			if strings.Contains(stdoutLine, "[build-stderr]") || strings.Contains(stdoutLine, "[ERROR]") {
				stdoutLine = regexp.MustCompile(`^\[[\d-\s:]+\]\s*`).ReplaceAllString(stdoutLine, "")
				outError += stdoutLine + "\n"
			}
		}
		panic("Database build failed:" + outError + stderr)
	}
	//writeTaskLog(&task, "Database build log:"+stderr)
	WriteTaskLog(task, "Database build completed:"+databasePath)
	return databasePath
}
