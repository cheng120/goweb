package main

import (
	"goweb/app/console"
	"goweb/app/http"
	"goweb/framework"
	"goweb/framework/provider/app"
	"goweb/framework/provider/config"
	"goweb/framework/provider/distributed"
	"goweb/framework/provider/env"
	"goweb/framework/provider/id"
	"goweb/framework/provider/kernel"
	"goweb/framework/provider/log"
	"goweb/framework/provider/trace"
)

func main() {
	// 初始化服务容器
	container := framework.NewContainer()
	// 绑定App服务提供者
	container.Bind(&app.AppProvider{})
	container.Bind(&env.EnvProvider{})
	container.Bind(&distributed.LocalDistributedProvider{})
	// // 后续初始化需要绑定的服务提供者...
	container.Bind(&config.ConfigProvider{})
	container.Bind(&id.IDProvider{})
	container.Bind(&trace.TraceProvider{})
	container.Bind(&log.LogServiceProvider{})
	// container.Bind(&orm.GormProvider{})
	// container.Bind(&redis.RedisProvider{})
	// container.Bind(&cache.CacheProvider{})
	// container.Bind(&ssh.SSHProvider{})


	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.KernelProvider{HttpEngine: engine})
	}

	// 运行root命令
	console.RunCommand(container)
}