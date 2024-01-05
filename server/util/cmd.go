package util

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func RunCmd(name string, arg ...string) (string, string, error) {
	cmd := exec.Command(name, arg...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Println("RunCmd: " + cmd.String())
	err := cmd.Run()

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	return strings.TrimSpace(outStr), strings.TrimSpace(errStr), err
}

var envMap map[string]string

func SetEnv(envToSet string) {
	//储存初始环境变量
	if envMap == nil {
		envMap = make(map[string]string)
		envs := os.Environ()
		for _, env := range envs {
			keyValue := strings.SplitN(env, "=", 2)
			if len(keyValue) == 2 {
				envMap[keyValue[0]] = keyValue[1]
			}
		}
	}

	//清空当前环境变量，为重新设置做准备
	{
		envs := os.Environ()
		for _, env := range envs {
			keyValue := strings.SplitN(env, "=", 2)
			if len(keyValue) == 2 { //reg := regexp.MustCompile(`^[A-Za-z0-9_.\s]+$`)
				err := os.Unsetenv(keyValue[0])
				if err != nil {
					log.Println("Error: " + err.Error() + " Code: os.Unsetenv" + keyValue[0])
				}
			}
		}
	}

	//还原系统原有环境变量
	for key, value := range envMap {
		err := os.Setenv(key, value)
		if err != nil {
			log.Println("Error: " + err.Error() + " Code: os.Setenv" + key + ", " + value)
		}
	}

	//设置新的变量
	envToSetList := strings.Split(envToSet, "\n")
	for _, envLine := range envToSetList {
		envLine = strings.TrimSpace(envLine)
		if !strings.HasPrefix(envLine, "#") {
			keyValue := strings.SplitN(envLine, "=", 2)
			if len(keyValue) == 2 {
				key := strings.TrimSpace(keyValue[0])
				value := strings.TrimSpace(keyValue[1])

				re := regexp.MustCompile(`\${(.+?)}`)
				value = re.ReplaceAllStringFunc(value, func(match string) string {
					return os.Getenv(match[2 : len(match)-1])
				})

				err := os.Setenv(key, value)
				if err != nil {
					log.Println("Error: " + err.Error() + " Code: os.Setenv" + key + ", " + value)
				}
			}
		}
	}
}
