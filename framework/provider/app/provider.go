package app

import (
	"goweb/framework"

	"goweb/framework/contract"
)

// AppProvider 提供App的具体实现方法
type AppProvider struct {
	BaseFolder string
}

// Register 注册HadeApp方法
func (h *AppProvider) Register(container framework.Container) framework.NewInstance {
	return NewApp
}

// Boot 启动调用
func (h *AppProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (h *AppProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (h *AppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, h.BaseFolder}
}

// Name 获取字符串凭证
func (h *AppProvider) Name() string {
	return contract.AppKey
}