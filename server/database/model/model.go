package model

import (
	"gorm.io/gorm"
	"time"
)

type Project struct {
	gorm.Model
	Url      string
	Owner    string
	Repo     string
	Mode     int //0监控release,1监控database
	Language string
	Command  string //编译命令
	Suite    StrArr
	Pause    bool

	LastAnalyzeTime           time.Time
	LastAnalyzeReleaseTag     string
	LastAnalyzeDatabaseCommit string

	LatestVersion           string
	LatestVersionErrorInfo  string //如果版本获取失败，用来保存失败信息
	LatestVersionUpdateTime time.Time
	LatestVersionCheckTime  time.Time
	LatestVersionCheckMode  int
}

type Setting struct {
	ID uint `gorm:"primarykey"`

	CodeQLCli      string
	CodeQLLib      string
	CodeQLPacks    string
	CodeQLSuite    string
	CodeQLDatabase string
	CodeQLResult   string
	EnvStr         string

	SystemToken string
	GithubToken string

	UpdateDetectionInterval float32 //项目更新检测最小间隔时间
	SkipVerifyTLS           bool

	AutoRecoveryTask  bool
	FirstReleaseCount int //第一次扫描release时，至多扫描多少个版本
	CronTaskSpec      string
	CronTaskNextTime  time.Time `gorm:"-"`

	AutoReadEmptyTask      bool
	AutoReadNoResultTask   bool
	AutoReadCompletedTask  bool
	AutoReadNoResultResult bool

	CodeQLAnalyzeOptions string
}

type Task struct {
	gorm.Model
	Status int    //0-队列中，1-进行中，2-完成，-1-错误
	Stage  int    //0-判断有无新版本，1-下载新版本，2-编译数据库，3-分析
	Logs   string //概要日志
	Manual bool   //是否手动触发的

	ProjectID       uint
	ProjectOwner    string
	ProjectRepo     string
	ProjectName     string
	ProjectLanguage string
	ProjectMode     int //0监控release,1监控database,2自定义任务
	ProjectCommand  string
	ProjectSuite    StrArr

	DatabasePath string //2需要的属性

	Versions         StrArr
	AnalyzedVersions StrArr

	TotalResultsCount int //本次任务总共产生的漏洞数量

	IsRead bool
}

type TaskResult struct {
	gorm.Model
	Version     string
	Commit      string
	FileName    string
	ResultCount int
	IsRead      bool

	TaskId uint
	Task   Task

	CodeQLSarif interface{} `gorm:"-"`
	FilePath    string      `gorm:"-"`
}

type Status struct {
	ID uint `gorm:"primarykey"`

	TotalTasks     int64
	CompletedTasks int64
	FailedTasks    int64
	TotalResults   int64
}
