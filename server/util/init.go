package util

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var storageDir string

func GetStorageDir() string {
	return storageDir
}

func InitStorageDir(path string) string {
	absPath, err := GetAbsolutePath(path)
	if err != nil {
		panic(err.Error())
	}

	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(absPath, os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	}

	storageDir = absPath
	return absPath
}

func InitSetting(storageDir string) {
	var settingCount int64
	database.DB.Model(&model.Setting{}).Count(&settingCount)
	if settingCount == 0 {
		setting := model.Setting{}
		setting.CodeQLCli = "codeql"
		setting.CodeQLPacks = filepath.Join(storageDir, "packs") + string(os.PathSeparator)
		setting.CodeQLSuite = filepath.Join(storageDir, "suites") + string(os.PathSeparator)
		setting.CodeQLDatabase = filepath.Join(storageDir, "databases") + string(os.PathSeparator)
		setting.CodeQLResult = filepath.Join(storageDir, "results") + string(os.PathSeparator)
		setting.EnvStr = `#MAVEN_HOME=/usr/local/maven
#PATH=${MAVEN_HOME}/bin:${PATH}

#MAVEN_HOME=D:\maven
#Path=${MAVEN_HOME}\bin:${Path}

#proxy=http://127.0.0.1:7890
http_proxy=${proxy}
https_proxy=${proxy}
`
		setting.SystemToken = ""
		setting.UpdateDetectionInterval = 60
		setting.SkipVerifyTLS = false
		setting.AutoRecoveryTask = false
		setting.FirstReleaseCount = 1
		setting.CronTaskSpec = "@weekly"
		setting.AutoReadEmptyTask = true
		setting.AutoReadNoResultTask = false
		setting.AutoReadCompletedTask = false
		setting.AutoReadNoResultResult = true
		setting.CodeQLAnalyzeOptions = "--ram 10240 --threads 0 --no-group-results"
		database.DB.Create(&setting)
	}
}

func InitStatus() {
	var statusCount int64
	database.DB.Model(&model.Status{}).Count(&statusCount)
	if statusCount == 0 {
		status := model.Status{}
		database.DB.Create(&status)
	}
}

func MakeDataDir(path string) {
	err := os.Mkdir(filepath.Join(path, "packs"), os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	err = os.Mkdir(filepath.Join(path, "suites"), os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	err = os.Mkdir(filepath.Join(path, "databases"), os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	err = os.Mkdir(filepath.Join(path, "results"), os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
}

func InitExample(dir string, examples *embed.FS) {
	if err := fs.WalkDir(examples, "examples", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		fileBytes, err := examples.ReadFile(path)
		if err != nil {
			panic(err.Error())
		}
		targetPath := filepath.Join(dir, path[strings.Index(path, "/")+1:])
		err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
		err = os.WriteFile(targetPath, fileBytes, 0644)
		if err != nil {
			panic(err.Error())
		}
		return err
	}); err != nil {
		panic(err.Error())
	}
}

func InitEnv() {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}
	SetEnv(setting.EnvStr)
}
