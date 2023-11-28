package processor

import (
	"choccy/server/database/model"
	"choccy/server/util"
)

// CheckDatabaseUpdates returns the new version name and download address
func CheckDatabaseUpdates(task *model.Task, project *model.Project) (string, string) {
	codeQLDatabase, err := util.GetGithubDatabase(task.ProjectOwner, task.ProjectRepo, task.ProjectLanguage)
	if err != nil {
		panic("Failed to get database:" + err.Error())
	}
	SetProjectLatestVersion(project, codeQLDatabase.CommitOid, codeQLDatabase.CreatedAt)
	return codeQLDatabase.CommitOid, codeQLDatabase.Url
}

func DownloadDatabase(task *model.Task, url string, commit string, databaseName string) string {
	WriteTaskLog(task, "Download Version:"+commit)
	databasePath, err := util.DownloadGithubDatabase(url, databaseName)
	if err != nil {
		panic("Database download failed:" + err.Error())
	}
	WriteTaskLog(task, "Download successfully, path:"+databasePath)
	return databasePath
}
