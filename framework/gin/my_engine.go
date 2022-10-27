package gin

import "goweb/framework"

func (engine *Engine) SetContainer(container framework.Container) {
	engine.Container = container
}

// engine实现container的绑定封装
func (engine *Engine) Bind(provider framework.ServiceProvider) error {
	return engine.Container.Bind(provider)
}

// IsBind 关键字凭证是否已经绑定服务提供者
func (engine *Engine) IsBind(key string) bool {
	return engine.Container.IsBind(key)
}