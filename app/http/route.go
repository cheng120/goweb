package http

import (
	"goweb/app/http/module/demo"

	"goweb/framework/gin"
)

func Routes(r *gin.Engine) {

	r.Static("/dist/", "./dist/")

	demo.Register(r)
}