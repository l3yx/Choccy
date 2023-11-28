package util

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func githubRequest(api string) ([]byte, error) {
	log.Println("GithubRequest: " + api)

	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	httpClient, err := GetHttpClient(time.Second * 10)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	if strings.TrimSpace(setting.GithubToken) != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", setting.GithubToken))
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(respBytes))
	}

	return respBytes, err
}

type GithubRelease struct {
	ID              int       `json:"id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	TarballUrl      string    `json:"tarball_url"`
	ZipballUrl      string    `json:"zipball_url"`
}

func GetGithubReleases(owner string, repo string) ([]GithubRelease, error) {
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?per_page=10", owner, repo)
	res, err := githubRequest(api)
	if err != nil {
		return nil, err
	}
	var githubReleases []GithubRelease
	err = json.Unmarshal(res, &githubReleases)
	if err != nil {
		return nil, err
	}
	return githubReleases, err
}

func GetGithubReleaseLatest(owner string, repo string) (*GithubRelease, error) {
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	res, err := githubRequest(api)
	if err != nil {
		return nil, err
	}
	var githubRelease GithubRelease
	err = json.Unmarshal(res, &githubRelease)
	if err != nil {
		return nil, err
	}
	return &githubRelease, nil
}

type CodeQlDatabase struct {
	ID          int       `json:"id"`
	Language    string    `json:"language"`
	ContentType string    `json:"content_type"`
	Size        int       `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Url         string    `json:"url"`
	CommitOid   string    `json:"commit_oid"`
}

func GetGithubDatabase(owner string, repo string, language string) (*CodeQlDatabase, error) {
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/code-scanning/codeql/databases/%s", owner, repo, language)
	res, err := githubRequest(api)
	if err != nil {
		return nil, err
	}

	var codeQlDatabase CodeQlDatabase
	err = json.Unmarshal(res, &codeQlDatabase)
	if err != nil {
		return nil, err
	}
	if codeQlDatabase.CommitOid == "" {
		codeQlDatabase.CommitOid = "null"
	}

	return &codeQlDatabase, err
}

func GetGithubDatabases(owner string, repo string) ([]CodeQlDatabase, error) {
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/code-scanning/codeql/databases", owner, repo)
	res, err := githubRequest(api)
	if err != nil {
		return nil, err
	}
	var codeQlDatabases []CodeQlDatabase
	err = json.Unmarshal(res, &codeQlDatabases)
	if err != nil {
		return nil, err
	}

	return codeQlDatabases, err
}

func DownloadGithubDatabase(databaseUrl, databaseName string) (string, error) {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	headers := make(map[string]string)
	headers["Accept"] = "application/zip"
	headers["X-GitHub-Api-Version"] = "2022-11-28"
	if strings.TrimSpace(setting.GithubToken) != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", setting.GithubToken)
	}

	downloadPath, err := downloadFile(databaseUrl, databaseName+".zip", headers)
	if err != nil {
		return "", err
	}
	defer os.Remove(downloadPath)
	settingPath, err := GetSettingPath()
	if err != nil {
		return "", err
	}
	databasePath := filepath.Join(settingPath.CodeQLDatabase, databaseName)
	_, err = os.Stat(databasePath)
	if err == nil {
		return "", fmt.Errorf("The database directory already exists:" + databasePath)
	}
	err = Unzip(downloadPath, databasePath, 1)
	if err != nil {
		isText, content := IsMiniTextFile(downloadPath)
		if isText {
			return "", fmt.Errorf("Failed to unzip file," + err.Error() + ". Detected that the file content is plain text, the content is:" + content)
		} else {
			return "", fmt.Errorf("Failed to unzip file," + err.Error())
		}
	}
	return databasePath, nil
}

func DownloadGithubTag(owner string, repo string, tagName string) (string, error) {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	headers := make(map[string]string)
	headers["X-GitHub-Api-Version"] = "2022-11-28"
	if strings.TrimSpace(setting.GithubToken) != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", setting.GithubToken)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/%s", owner, repo, tagName)
	codeName := fmt.Sprintf("%s__%s__%s", owner, repo, MakeValidFilename(tagName))
	downloadPath, err := downloadFile(url, codeName+".zip", headers)
	if err != nil {
		return "", err
	}
	defer os.Remove(downloadPath)

	codePath := filepath.Join(filepath.Dir(downloadPath), codeName)
	err = Unzip(downloadPath, codePath, 1)
	if err != nil {
		isText, content := IsMiniTextFile(downloadPath)
		if isText {
			return "", fmt.Errorf("Failed to unzip file," + err.Error() + ". Detected that the file content is plain text, the content is:" + content)
		} else {
			return "", fmt.Errorf("Failed to unzip file," + err.Error())
		}
	}

	return codePath, nil
}

type GithubTagCommit struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type GithubTag struct {
	Name       string          `json:"name"`
	ZipballUrl string          `json:"zipball_url"`
	TarballUrl string          `json:"tarball_url"`
	Commit     GithubTagCommit `json:"commit"`
	NodeId     string          `json:"node_id"`
}

func GetGithubTag(owner string, repo string, tagName string) (*GithubTag, error) {
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo)
	res, err := githubRequest(api)
	if err != nil {
		return nil, err
	}
	var githubTags []GithubTag
	err = json.Unmarshal(res, &githubTags)
	if err != nil {
		return nil, err
	}

	var targetTag GithubTag
	for _, githubTag := range githubTags {
		if githubTag.Name == tagName {
			targetTag = githubTag
			break
		}
	}

	return &targetTag, err
}
