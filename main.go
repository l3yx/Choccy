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
	//Get command line parameters, only addr and token
	addr := flag.String("addr", "0.0.0.0:80", "Listening address and port")
	token := flag.String("token", "", "System Token")
	flag.Parse()

	//Initialize data directory, database, environment variables
	{
		storageDir := "choccy_data"
		util.InitStorageDir(storageDir)
		// isNew, is true only when choccy.db does not exist, that is, when the program is started for the first time
		isNew := database.InitDB(filepath.Join(util.GetStorageDir(), "choccy.db"))
		if isNew { // The program starts for the first time and performs various initializations.
			if strings.TrimSpace(*token) == "" { //If the user does not specify a token, it will be randomly generated
				*token = util.RandomString(16)
			}
			util.MakeDataDir(util.GetStorageDir())
			util.InitExample(util.GetStorageDir(), &examples)
		}
		util.InitSetting(storageDir)
		util.InitStatus()
		util.InitEnv()
	}

	//Set system token
	{
		if strings.TrimSpace(*token) != "" {
			database.DB.Model(model.Setting{}).
				Select("SystemToken").Where(true).Updates(model.Setting{SystemToken: *token})
		}
	}

	//Initialize the task executor and resume unexecuted tasks
	{
		taskmanager.InitTask()
		err := taskmanager.SetCronTask()
		if err != nil {
			log.Println("Error: " + err.Error())
		}
		go taskmanager.Consumer()
	}

	//Register route
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
