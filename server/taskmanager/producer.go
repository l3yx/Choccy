package taskmanager

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func AddTask(project *model.Project, manual bool) (bool, error) {
	filters := make(map[string]interface{})
	filters["Status"] = []int{0, 1}
	filters["project_id"] = project.ID
	filters["project_owner"] = project.Owner
	filters["project_repo"] = project.Repo
	filters["project_language"] = project.Language
	filters["project_mode"] = project.Mode
	filters["project_command"] = project.Command
	filters["project_suite"] = project.Suite
	var count int64
	result := database.DB.Model(&model.Task{}).Where(filters).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	if count > 0 {
		return false, nil
	}

	task := model.Task{
		Manual:           manual,
		ProjectID:        project.ID,
		ProjectOwner:     project.Owner,
		ProjectRepo:      project.Repo,
		ProjectName:      fmt.Sprintf("%s/%s", project.Owner, project.Repo),
		ProjectLanguage:  project.Language,
		ProjectMode:      project.Mode,
		ProjectCommand:   project.Command,
		ProjectSuite:     project.Suite,
		Versions:         []string{},
		AnalyzedVersions: []string{},
	}
	result = database.DB.Save(&task)
	if result.Error != nil {
		return false, result.Error
	}

	CH <- task.ID

	return true, nil
}

func AddCustomTask(databasePath string, suites []string, name string) (bool, error) {
	databaseYml := filepath.Join(databasePath, "codeql-database.yml")
	fileBytes, err := os.ReadFile(databaseYml)
	if err != nil {
		return false, err
	}
	var language string
	match := regexp.MustCompile(`primaryLanguage\s*:\s*(.+)`).FindSubmatch(fileBytes)
	if len(match) == 0 {
		language = "unknown"
	} else {
		language = strings.TrimSpace(string(match[1]))
	}

	task := model.Task{
		Manual:           true,
		ProjectName:      name,
		ProjectLanguage:  language,
		ProjectMode:      2,
		ProjectSuite:     suites,
		Versions:         []string{},
		AnalyzedVersions: []string{},

		DatabasePath: databasePath,
	}

	result := database.DB.Save(&task)
	if result.Error != nil {
		return false, result.Error
	}

	CH <- task.ID
	return true, nil
}
