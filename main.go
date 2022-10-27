package main

import (
	"goweb/app/console"
	"goweb/app/http"
	"goweb/framework"
	"goweb/framework/provider/app"
	"goweb/framework/provider/distributed"
	"goweb/framework/provider/env"
	"goweb/framework/provider/kernel"
)

func main() {
	// 初始化服务容器
	container := framework.NewContainer()
	// 绑定App服务提供者
	container.Bind(&app.HadeAppProvider{})
	container.Bind(&env.EnvProvider{})
	container.Bind(&distributed.LocalDistributedProvider{})
	// 后续初始化需要绑定的服务提供者...

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.KernelProvider{HttpEngine: engine})
	}

	// 运行root命令
	console.RunCommand(container)
}