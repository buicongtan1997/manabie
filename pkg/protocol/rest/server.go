package rest

import (
	apiV1 "github.com/buicongtan1997/manabie/pkg/api/v1/controllers"
	"github.com/buicongtan1997/manabie/pkg/logger"
	"github.com/buicongtan1997/manabie/pkg/protocol/rest/middleware"
	"github.com/buicongtan1997/manabie/pkg/services/v1"
	"github.com/gin-gonic/gin"
)

// RunServer runs HTTP/REST gateway
func RunServer() error {
	router := gin.New()

	//add middleware
	router.Use(middleware.AddRequestID())
	router.Use(middleware.AddLogger(logger.Log))

	api := router.Group("")
	apiV1.RegisterProductController(api, v1.NewProductService())
	return router.Run()
}
