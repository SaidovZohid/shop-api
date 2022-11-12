package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "golang/docs"
)

type handler struct {
	db *DBManager
}

// @title           Swagger for Shop api
// @version         1.0
// @description     This is a Shop service Api
// @host      localhost:8080
func NewServer(db *DBManager) *gin.Engine {
	r := gin.Default()

	h := handler{
		db: db,
	}
	r.GET("/shop/customer/:id", h.GetCustomer)
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
