package router

import (
	"choccy/server/handler"
	"choccy/server/middleware"
	"github.com/gin-gonic/gin"
)

func SetupApiRoutes(r *gin.Engine) {
	apiRouter := r.Group("/api", middleware.AuthMiddleware)
	{
		apiRouter.GET("/project", handler.GetProjects)
		apiRouter.POST("/project", handler.SaveProject)
		apiRouter.DELETE("/project", handler.DeleteProject)

		apiRouter.GET("/setting", handler.GetSetting)
		apiRouter.POST("/setting", handler.SaveSetting)
		apiRouter.POST("/setting/test", handler.TestSetting)

		apiRouter.GET("/query", handler.GetQueries)
		apiRouter.GET("/query/content", handler.GetQueryContent)

		apiRouter.GET("/suite", handler.GetSuites)
		apiRouter.DELETE("/suite", handler.DeleteSuite)
		apiRouter.POST("/suite", handler.CreateSuite)
		apiRouter.POST("/suite/rename", handler.RenameSuite)
		apiRouter.GET("/suite/content", handler.GetSuiteContent)
		apiRouter.POST("/suite/content", handler.SaveSuiteContent)
		apiRouter.GET("/suite/resolve", handler.ResolveSuiteQueries)

		apiRouter.GET("/database", handler.GetDatabases)

		apiRouter.GET("/result", handler.GetResult)
		apiRouter.GET("/result/unread", handler.GetResultUnread)
		apiRouter.DELETE("/result", handler.DeleteResult)
		apiRouter.POST("/result/read", handler.SetResultIsRead)

		apiRouter.GET("/task/run", handler.RunTask)
		apiRouter.GET("/task", handler.GetTasks)
		apiRouter.GET("/task/unread", handler.GetTaskUnread)
		apiRouter.POST("/task/read", handler.SetTaskIsRead)

		apiRouter.GET("/notifications", handler.GetNotifications)
	}
}
