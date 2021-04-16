package server

import (
	"github.com/gin-gonic/gin"
)

func Init() {
	e := gin.Default()
	gConfig := e.Group("/task/back")
	{
		gConfig.POST("/register", register)
		gConfig.POST("/feedback", feedBack)
		gConfig.POST("/start", start)
		gConfig.POST("/end", stop)
		gConfig.POST("/update", update)

		gConfig.GET("/health", health)
	}
	issue := e.Group("/issue")
	{
		issue.POST("/handle", IssueHandle)
	}
	e.Run()
}
