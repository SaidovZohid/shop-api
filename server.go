package main

import (
	_ "golang/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r.POST("/shop/customer/create", h.CreateCustomer)
	r.PUT("/shop/customer/update/:id", h.UpdateCustomer)
	r.DELETE("/shop/customer/delete/:id", h.DeleteCustomer)
	r.GET("/shop/customer/all", h.GetAllCustomer)

	r.POST("/shop/category/create", h.CreateCategory)
	r.GET("/shop/category/:id", h.GetCategory)
	r.PUT("/shop/category/update/:id", h.UpdateCategory)

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
