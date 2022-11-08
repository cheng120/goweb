package http

import (
	"goweb/app/http/module/demo"

	"goweb/framework/gin"
	"goweb/framework/middleware/static"
)

func Routes(r *gin.Engine) {

	// /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))
	// r.Static("/dist/", "./dist/")

	demo.Register(r)
}