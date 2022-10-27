package kernel

import (
	"goweb/framework/gin"
	"net/http"
)

// 引擎服务
type KernelService struct{
	Engine *gin.Engine
}

// 初始化 web 引擎服务实例
func NewKernelService(params ...interface{}) (interface{},error) {
	httpEngine := params[0].(*gin.Engine)
	return &KernelService{Engine: httpEngine}, nil
}

// 返回 web 引擎
func (ks *KernelService) HttpEngine() http.Handler{
	return ks.Engine
}