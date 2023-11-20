package util

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func downloadFile(downloadUrl string, fileName string, headers map[string]string) (string, error) {
	log.Println("Download: " + downloadUrl)
	tmpDir, err := GetTmpDir()
	if err != nil {
		return "", err
	}

	httpClient, err := GetHttpClient(time.Minute * 30)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return "", err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(tmpDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}
	defer file.Close()
	defer resp.Body.Close()

	return filePath, nil
}
