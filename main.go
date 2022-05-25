package main

import (
	"github.com/labstack/echo/v4"
	"github.com/raa11dev/crud-echo/database"
	"github.com/raa11dev/crud-echo/routes"
)

func main() {
	e := echo.New()

	//connect to mongoDB
	database.ConnectDB()

	//routes
	routes.ProductRoute(e)

	e.Logger.Fatal(e.Start(":8000"))
}
