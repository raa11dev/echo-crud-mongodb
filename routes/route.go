package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/raa11dev/crud-echo/controllers"
)

func ProductRoute(e *echo.Echo) {
	e.POST("/product", controllers.CreateProduct)
	e.PUT("/product/:id", controllers.UpdateProduct)
	e.DELETE("/product/:id", controllers.DeleteProduct)
	e.GET("/product/:id", controllers.GetProduct)
}
