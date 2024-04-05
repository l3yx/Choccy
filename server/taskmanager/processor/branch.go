package processor

import (
	"choccy/server/database/model"
	"choccy/server/util"
)

// CheckBranchUpdates 返回新版本名，
func CheckDefaultBranchUpdates(task *model.Task, project *model.Project) string {
	branch, err := util.GetGithubDefaultBranch(task.ProjectOwner, task.ProjectRepo)
	if err != nil {
		panic("获取数据库失败：" + err.Error())
	}

	SetProjectLatestVersion(project, branch.Commit.Sha, branch.Commit.Commit.Committer.Date)
	return branch.Commit.Sha
}

func DownloadCommit(task *model.Task, commit string) string {
	WriteTaskLog(task, "下载版本："+commit)
	tagSourcePath, err := util.DownloadGithubTag(task.ProjectOwner, task.ProjectRepo, commit)
	if err != nil {
		panic("下载失败：" + err.Error())
	}
	WriteTaskLog(task, "下载成功，路径："+tagSourcePath)
	return tagSourcePath
}
