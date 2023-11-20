package util

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetAbsolutePath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	var absolutePath string
	if filepath.IsAbs(path) {
		absolutePath = path
	} else {
		absolutePath = filepath.Join(filepath.Dir(executablePath), path)
	}

	absPath, err := filepath.Abs(absolutePath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

type SettingPath struct {
	CodeQLLib      string
	CodeQLPacks    string
	CodeQLSuite    string
	CodeQLDatabase string
	CodeQLCli      string
	CodeQlResult   string
}

func GetSettingPath() (*SettingPath, error) {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		return nil, result.Error
	}
	codeQLPacks, err := GetAbsolutePath(setting.CodeQLPacks)
	if err != nil {
		return nil, err
	}
	codeQLLib := ""
	if strings.TrimSpace(setting.CodeQLLib) != "" {
		codeQLLib, err = GetAbsolutePath(setting.CodeQLLib)
		if err != nil {
			return nil, err
		}
	}
	codeQLSuite, err := GetAbsolutePath(setting.CodeQLSuite)
	if err != nil {
		return nil, err
	}
	codeQLDatabase, err := GetAbsolutePath(setting.CodeQLDatabase)
	if err != nil {
		return nil, err
	}
	codeQLResult, err := GetAbsolutePath(setting.CodeQLResult)
	if err != nil {
		return nil, err
	}
	codeQLCli, err := GetAbsolutePath(setting.CodeQLCli)
	if err != nil {
		return nil, err
	}
	if setting.CodeQLCli == "codeql" {
		codeQLCli = "codeql"
	}
	settingPath := &SettingPath{CodeQLLib: codeQLLib, CodeQLPacks: codeQLPacks, CodeQLSuite: codeQLSuite, CodeQLDatabase: codeQLDatabase, CodeQLCli: codeQLCli, CodeQlResult: codeQLResult}
	return settingPath, nil
}

func GetCodeQL(path string) (string, error) {
	if path == "codeql" {
		return "codeql", nil
	}
	return GetAbsolutePath(path)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type FileInfo struct {
	Name    string
	Path    string
	ModTime time.Time
	Extra   map[string]interface{}
}

func ListFiles(folderOnly bool, fileOnly bool, allowedExtensions []string, path string, sortBy string, sortOrder string, pageSize int, pageNumber int) ([]FileInfo, int, error) {
	var items []FileInfo

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, 0, err
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, 0, err
	}

	for _, file := range files {
		if folderOnly {
			if !file.IsDir() {
				continue
			}
		}
		if fileOnly {
			if file.IsDir() {
				continue
			}
		}

		if !file.IsDir() && !contains(allowedExtensions, "*") {
			ext := filepath.Ext(file.Name())
			if !contains(allowedExtensions, ext) {
				continue
			}
		}

		extra := make(map[string]interface{})
		if file.IsDir() {
			databaseYml := filepath.Join(path, file.Name(), "codeql-database.yml")
			if _, err := os.Stat(databaseYml); os.IsNotExist(err) {
				continue
			}
			fileBytes, err := os.ReadFile(databaseYml)
			if err != nil {
				return nil, 0, err
			}
			match := regexp.MustCompile(`primaryLanguage\s*:\s*(.+)`).FindSubmatch(fileBytes)
			if len(match) == 0 {
				extra["database_language"] = "unknown"
			} else {
				extra["database_language"] = strings.TrimSpace(string(match[1]))
			}
			match = regexp.MustCompile(`baselineLinesOfCode\s*:\s*(.+)`).FindSubmatch(fileBytes)
			if len(match) == 0 {
				extra["database_linesOfCode"] = "unknown"
			} else {
				extra["database_linesOfCode"], _ = strconv.Atoi(strings.TrimSpace(string(match[1])))
			}
			match = regexp.MustCompile(`cliVersion\s*:\s*(.+)`).FindSubmatch(fileBytes)
			if len(match) == 0 {
				extra["database_cliVersion"] = "unknown"
			} else {
				extra["database_cliVersion"] = strings.TrimSpace(string(match[1]))
			}
			match = regexp.MustCompile(`finalised\s*:\s*(.+)`).FindSubmatch(fileBytes)
			if len(match) == 0 {
				extra["database_finalised"] = "unknown"
			} else {
				extra["database_finalised"] = strings.TrimSpace(string(match[1]))
			}
		} else if filepath.Ext(file.Name()) == ".qls" {
			qlsFile := filepath.Join(path, file.Name())
			fileBytes, err := os.ReadFile(qlsFile)
			if err != nil {
				return nil, 0, err
			}
			re := regexp.MustCompile(`description\s*:\s*(.+)`)
			match := re.FindSubmatch(fileBytes)
			if len(match) == 0 {
				extra["suite_description"] = "unknown"
			} else {
				extra["suite_description"] = string(match[1])
			}
		}

		info, err := file.Info()
		if err != nil {
			return nil, 0, err
		}

		items = append(items, FileInfo{Name: file.Name(), ModTime: info.ModTime(), Path: filepath.Join(path, file.Name()), Extra: extra})
	}

	descending := false
	if sortOrder == "descending" {
		descending = true
	}
	sort.Slice(items, func(i, j int) bool {
		switch sortBy {
		case "Name":
			res := items[i].Name < items[j].Name
			if descending {
				return !res
			} else {
				return res
			}
		case "Extra.suite_description":
			res := items[i].Extra["suite_description"].(string) < items[j].Extra["suite_description"].(string)
			if descending {
				return !res
			} else {
				return res
			}
		case "Extra.database_language":
			res := items[i].Extra["database_language"].(string) < items[j].Extra["database_language"].(string)
			if descending {
				return !res
			} else {
				return res
			}
		case "Extra.database_linesOfCode":
			iLines, _ := items[i].Extra["database_linesOfCode"].(int)
			jLines, _ := items[j].Extra["database_linesOfCode"].(int)
			res := iLines < jLines
			if descending {
				return !res
			} else {
				return res
			}
		case "Extra.database_cliVersion":
			res := items[i].Extra["database_cliVersion"].(string) < items[j].Extra["database_cliVersion"].(string)
			if descending {
				return !res
			} else {
				return res
			}
		case "Extra.database_finalised":
			res := items[i].Extra["database_finalised"].(string) < items[j].Extra["database_finalised"].(string)
			if descending {
				return !res
			} else {
				return res
			}
		default:
			res := items[i].ModTime.Before(items[j].ModTime)
			if descending {
				return !res
			} else {
				return res
			}
		}
	})

	totalFolders := len(items)

	if pageSize != -1 {
		startIndex := (pageNumber - 1) * pageSize
		endIndex := startIndex + pageSize

		if endIndex > totalFolders {
			endIndex = totalFolders
		}

		if startIndex < 0 || startIndex >= totalFolders {
			return items, totalFolders, nil
		}

		return items[startIndex:endIndex], totalFolders, nil
	}
	return items, totalFolders, nil
}

func GetTmpDir() (string, error) {
	tmpDir := filepath.Join(GetStorageDir(), "tmp")

	_, err := os.Stat(tmpDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(tmpDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return tmpDir, nil
}

func MakeValidFilename(filename string) string {
	re := regexp.MustCompile(`[\s<>:"/\\|?*]`)
	validFilename := re.ReplaceAllString(filename, "_")
	return validFilename
}
