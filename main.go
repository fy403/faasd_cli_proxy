package main

import (
	"proxy/api"
	"proxy/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	auth := middleware.CookieAuth(middleware.Accounts(map[string]string{
		"admin": "123456",
	}))

	app := gin.New()
	app.Use(gin.Recovery(), gin.Logger())

	_api := app.Group("/api")
	_api.Use(auth)
	// 设置路由
	// /api/action
	_api.GET("/login", func(c *gin.Context) {})
	_api.GET("/logout", func(c *gin.Context) {})

	testHandler := new(api.TestHandler)
	_api.GET("/test", api.WarpHandle(testHandler.Info))

	// URI: /api/faasd/acton
	faasd := _api.Group("/faasd")
	faasCliHandler := new(api.FassCliHandler)
	faasd.POST("/new", api.WarpHandle(faasCliHandler.New))
	faasd.POST("/write", api.WarpHandle(faasCliHandler.Write))
	faasd.POST("/build", api.WarpHandle(faasCliHandler.Build))
	faasd.POST("/push", api.WarpHandle(faasCliHandler.Push))
	faasd.POST("/deploy", api.WarpHandle(faasCliHandler.Deploy))
	faasd.GET("/list", api.WarpHandle(faasCliHandler.GetAllInvokeInfo))
	faasd.GET("/describe", api.WarpHandle(faasCliHandler.GetInvokeInfo))
	faasd.POST("/up", api.WarpHandle(faasCliHandler.Up))
	faasd.GET("/support", api.WarpHandle(faasCliHandler.SupportedLang))
	app.Run(":8080")
}
