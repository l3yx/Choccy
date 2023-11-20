package main

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"choccy/server/middleware"
	"choccy/server/router"
	"choccy/server/taskmanager"
	"choccy/server/util"
	"embed"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"path/filepath"
	"strings"
)

//go:embed assets/*
var assets embed.FS

//go:embed index.html favicon.ico
var index embed.FS

//go:embed examples/*
var examples embed.FS

func run() {
	//获取命令行参数，只有addr和token
	addr := flag.String("addr", "0.0.0.0:80", "监听地址和端口")
	token := flag.String("token", "", "系统Token")
	flag.Parse()

	//初始化数据目录，数据库，环境变量
	{
		storageDir := "choccy_data"
		util.InitStorageDir(storageDir)
		// isNew，只在不存在choccy.db时为true，也就是程序第一次启动的时候
		isNew := database.InitDB(filepath.Join(util.GetStorageDir(), "choccy.db"))
		if isNew { // 程序第一次启动进行各项初始化
			if strings.TrimSpace(*token) == "" { //如果用户不指定token，则随机生成
				*token = util.RandomString(16)
			}
			util.MakeDataDir(util.GetStorageDir())
			util.InitExample(util.GetStorageDir(), &examples)
		}
		util.InitSetting(storageDir)
		util.InitStatus()
		util.InitEnv()
	}

	//设置系统Token
	{
		if strings.TrimSpace(*token) != "" {
			database.DB.Model(model.Setting{}).
				Select("SystemToken").Where(true).Updates(model.Setting{SystemToken: *token})
		}
	}

	//初始化任务执行器，恢复未执行的任务
	{
		taskmanager.InitTask()
		err := taskmanager.SetCronTask()
		if err != nil {
			log.Println("Error: " + err.Error())
		}
		go taskmanager.Consumer()
	}

	//注册路由
	{
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		r.Use(middleware.ErrorMiddleware)
		r.Use(middleware.GetCorsMiddleware())

		router.SetupFrontendRoutes(r, &assets, &index)
		router.SetupApiRoutes(r)

		var setting model.Setting
		result := database.DB.Take(&setting)
		if result.Error != nil {
			panic(result.Error.Error())
		}

		log.Println("Token : " + setting.SystemToken)
		log.Println("Listen: " + *addr)
		err := r.Run(*addr)
		log.Fatalln(err)
	}

}

func main() {
	run()
}
