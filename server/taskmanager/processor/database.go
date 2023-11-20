package processor

import (
	"choccy/server/database/model"
	"choccy/server/util"
)

// CheckDatabaseUpdates 返回新版本名，下载地址
func CheckDatabaseUpdates(task *model.Task, project *model.Project) (string, string) {
	codeQLDatabase, err := util.GetGithubDatabase(task.ProjectOwner, task.ProjectRepo, task.ProjectLanguage)
	if err != nil {
		panic("获取数据库失败：" + err.Error())
	}
	SetProjectLatestVersion(project, codeQLDatabase.CommitOid, codeQLDatabase.CreatedAt)
	return codeQLDatabase.CommitOid, codeQLDatabase.Url
}

func DownloadDatabase(task *model.Task, url string, commit string, databaseName string) string {
	WriteTaskLog(task, "下载版本："+commit)
	databasePath, err := util.DownloadGithubDatabase(url, databaseName)
	if err != nil {
		panic("数据库下载失败：" + err.Error())
	}
	WriteTaskLog(task, "下载成功，路径："+databasePath)
	return databasePath
}
