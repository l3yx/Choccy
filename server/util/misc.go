package util

import (
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[r.Intn(len(charset))])
	}
	return sb.String()
}

func Reverse(arr *[]string) {
	var temp string
	length := len(*arr)
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}

func IsCodeQLDatabaseExists(databaseName string) string {
	settingPath, err := GetSettingPath()
	if err != nil {
		panic(err.Error())
	}
	databasePath := filepath.Join(settingPath.CodeQLDatabase, databaseName)
	_, err = os.Stat(databasePath)
	if os.IsNotExist(err) {
		return ""
	} else if err != nil {
		panic(err.Error())
	}
	return databasePath
}

func IsMiniTextFile(path string) (bool, string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, ""
	}
	fileSize := fileInfo.Size()
	if fileSize > 1024 {
		return false, ""
	}
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return false, ""
	}
	if regexp.MustCompile("[[:^ascii:]]").Match(fileContent) {
		return false, ""
	}
	return true, string(fileContent)
}
