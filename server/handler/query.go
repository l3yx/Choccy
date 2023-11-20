package handler

import (
	"choccy/server/util"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strings"
)

func GetQueryContent(c *gin.Context) {
	path := c.DefaultQuery("path", "")
	path = filepath.Clean(path)

	if filepath.Ext(path) != ".ql" {
		panic("Must be a ql file")
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"data": string(bytes),
	})
}

func GetQueries(c *gin.Context) {
	path := c.DefaultQuery("path", "")

	if path != "" {
		data := util.ResolvePackQueries(path)
		c.JSON(200, gin.H{
			"data": data,
		})
	} else {
		data := make(map[string]map[string][]string)
		packs := util.ResolvePacks("query")
		for key, paths := range packs {
			scope := "/"
			var pack string
			name := strings.SplitN(key, "/", 2)
			if len(name) == 2 {
				scope = name[0]
				pack = name[1]
			} else {
				pack = key
			}

			if _, ok := data[scope]; !ok {
				data[scope] = make(map[string][]string)
			}
			if _, ok := data[scope][pack]; !ok {
				data[scope][pack] = []string{}
			}
			data[scope][pack] = append(data[scope][pack], paths...)
		}

		c.JSON(200, gin.H{
			"data": data,
		})
	}

}
