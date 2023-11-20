package util

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func getCodeQLCli() string {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	codeQLCli, err := GetCodeQL(setting.CodeQLCli)
	if err != nil {
		panic(err.Error())
	}
	return codeQLCli
}

func getSearchPath() string {
	settingPath, err := GetSettingPath()
	if err != nil {
		panic(err.Error())
	}

	var searchPathList []string
	if strings.TrimSpace(settingPath.CodeQLLib) != "" {
		searchPathList = append(searchPathList, settingPath.CodeQLLib)
	}
	if strings.TrimSpace(settingPath.CodeQLPacks) != "" {
		searchPathList = append(searchPathList, strings.Split(settingPath.CodeQLPacks, "\n")...)
	}

	searchPath := strings.Join(searchPathList, string(os.PathListSeparator))

	return searchPath
}

func ResolvePacks(kind string) map[string][]string {
	args := []string{"resolve", "qlpacks",
		"--format", "json",
		"-q",
		"--kind", kind}

	searchPath := getSearchPath()
	if searchPath != "" {
		args = append(args, "--search-path", searchPath)
	}

	stdout, stderr, err := RunCmd(getCodeQLCli(), args...)
	if err != nil {
		panic(err.Error() + "\n" + stdout + "\n" + stderr)
	}
	if stderr != "" {
		panic(stderr)
	}

	var data map[string][]string
	err = json.Unmarshal([]byte(stdout), &data)
	if err != nil {
		panic(err.Error())
	}

	return data
}

func ResolvePackQueries(dir string) []string {
	args := []string{"resolve", "queries",
		dir,
		"--format", "json",
		"-q",
	}

	searchPath := getSearchPath()
	if searchPath != "" {
		args = append(args, "--search-path", searchPath)
	}

	stdout, stderr, err := RunCmd(getCodeQLCli(), args...)
	if err != nil {
		panic(err.Error() + "\n" + stdout + "\n" + stderr)
	}
	if stderr != "" {
		panic(stderr)
	}

	var data []string
	err = json.Unmarshal([]byte(stdout), &data)
	if err != nil {
		panic(err.Error())
	}

	return data
}

func ResolveSuiteQueries(path string) []string {
	args := []string{"resolve", "queries",
		path,
		"--format", "json",
		"-q",
	}

	searchPath := getSearchPath()
	if searchPath != "" {
		args = append(args, "--search-path", searchPath)
	}

	stdout, stderr, err := RunCmd(getCodeQLCli(), args...)
	if err != nil {
		panic(err.Error() + "\n" + stdout + "\n" + stderr)
	}
	if stderr != "" {
		panic(stderr)
	}

	var data []string
	err = json.Unmarshal([]byte(stdout), &data)
	if err != nil {
		panic(err.Error())
	}

	return data
}

func GetCodeQLVersionByPath(path string) string {
	stdout, stderr, err := RunCmd(path, "version", "-q")
	if err != nil {
		panic(err.Error() + "\n" + stdout + "\n" + stderr)
	}
	if stderr != "" {
		panic(stderr)
	}
	return stdout
}

func GetCodeQLVersion() string {
	stdout, stderr, err := RunCmd(getCodeQLCli(), "version", "-q")
	if err != nil {
		panic(err.Error() + "\n" + stdout + "\n" + stderr)
	}
	if stderr != "" {
		panic(stderr)
	}
	return stdout
}

func ResolveQueryMetadata(path string) (map[string]interface{}, error) {
	//https://github.com/github/codeql/blob/main/docs/query-metadata-style-guide.md#metadata-area
	metadata := make(map[string]interface{})

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return metadata, err
	}
	match := regexp.MustCompile(`/\*\*((?s).*?)\*/`).FindSubmatch(fileBytes)
	if len(match) == 0 {
		return metadata, nil
	}

	metadataString := string(match[1])
	metadataString = regexp.MustCompile(`(?m)^[\s*]*`).ReplaceAllString(metadataString, "")

	res := regexp.MustCompile(`@(\S+)\s+([^@]+)`).FindAllStringSubmatch(metadataString, -1)
	for _, match := range res {
		metadata[strings.TrimSpace(match[1])] = strings.TrimSpace(match[2])
	}

	metadata["path"] = path
	metadata["content"] = string(fileBytes)

	value, exists := metadata["description"]
	if exists {
		metadata["description"] = strings.ReplaceAll(value.(string), "\n", " ")
	}

	value, exists = metadata["tags"]
	if exists {
		metadata["tags"] = regexp.MustCompile(`\s+`).Split(value.(string), -1)
	}

	return metadata, nil
}

func DatabaseCreate(source, language, command, databaseName string) (string, string, error, string) {
	settingPath, err := GetSettingPath()
	if err != nil {
		return "", "", err, ""
	}

	databasePath := filepath.Join(settingPath.CodeQLDatabase, MakeValidFilename(databaseName))

	args := []string{"database", "create",
		databasePath,
		"--language", language,
		"--source-root", source,
		"--threads", "0",
	}

	searchPath := getSearchPath()
	if searchPath != "" {
		args = append(args, "--search-path", searchPath)
	}

	if strings.TrimSpace(command) != "" {
		args = append(args, "--command", command)
	}

	stdout, stderr, err := RunCmd(getCodeQLCli(), args...)

	return stdout, stderr, err, databasePath
}

func DatabaseAnalyze(databasePath, qls, outName string) (string, string, error, string) {
	settingPath, err := GetSettingPath()
	if err != nil {
		return "", "", err, ""
	}

	outputPath := filepath.Join(settingPath.CodeQlResult, MakeValidFilename(outName))

	args := []string{"database", "analyze",
		databasePath,
		qls,
		"--format", "sarif-latest",
		"--output", outputPath,
		"--threads", "0",
		"--rerun",
		"--sarif-add-snippets",
		"--no-group-results",
	}

	searchPath := getSearchPath()
	if searchPath != "" {
		args = append(args, "--search-path", searchPath)
	}

	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		return "", "", result.Error, ""
	}

	if strings.TrimSpace(setting.CodeQLAnalyzeOptions) != "" {
		args = append(args, strings.Split(strings.TrimSpace(setting.CodeQLAnalyzeOptions), " ")...)
	}

	stdout, stderr, err := RunCmd(getCodeQLCli(), args...)
	return stdout, stderr, err, outputPath
}

func GenerateSuite(suiteNames []string) (string, error) {
	settingPath, err := GetSettingPath()
	if err != nil {
		return "", err
	}

	tmpDir, err := GetTmpDir()
	if err != nil {
		return "", err
	}

	content := ""
	for _, suiteName := range suiteNames {
		content += "- import: " + filepath.Join(settingPath.CodeQLSuite, suiteName) + "\n"
	}

	suitePath := filepath.Join(tmpDir, RandomString(5)+".qls")
	file, err := os.Create(suitePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		return "", err
	}

	return suitePath, nil
}
