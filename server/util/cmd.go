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
	//Storage Initial Environment Variables
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

	//Empty the current environment variable in preparation for reset
	{
		envs := os.Environ()
		for _, env := range envs {
			keyValue := strings.SplitN(env, "=", 2)
			if len(keyValue) == 2 { //reg := regexp.MustCompile(`^[A-Za-z0-9_.\s]+$`)
				err := os.Unsetenv(keyValue[0])
				if err != nil {
					panic(err.Error())
				}
			}
		}
	}

	//restore system original environment variable
	for key, value := range envMap {
		err := os.Setenv(key, value)
		if err != nil {
			panic(err.Error())
		}
	}

	//Setting new variables
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
					panic(err.Error())
				}
			}
		}
	}
}
