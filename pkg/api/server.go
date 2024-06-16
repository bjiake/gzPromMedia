package api

import (
	"awesomeProject/pkg/api/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.Handler) *ServerHTTP {
	engine := gin.New()

	// Use logger from Gin
	engine.Use(gin.Logger())

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://172.20.196.70:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	engine.Use(cors.New(corsConfig))

	//Account methods
	engine.POST("/registration", userHandler.Registration)
	engine.POST("/login", userHandler.Login)

	// Use middleware from Gin
	engine.Use(userHandler.AuthMiddleware())

	engine.PUT("/account/:accountId", userHandler.PutAccount)
	engine.GET("/account/:accountId", userHandler.GetAccount)
	engine.DELETE("/account/:accountId", userHandler.DeleteAccount)
	engine.PUT("/subscribe:/accountId", userHandler.Subscribe)
	engine.PUT("/unSubscribe:/accountId", userHandler.UnSubscribe)
	engine.GET("/isBirthday", userHandler.CheckBirthDay)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run("127.0.0.1:8001")
}
