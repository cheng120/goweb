package gin

/// engine 实现容器 功能 end

// context 实现 container 的几个封装

func (ctx *Context) Make(key string) (interface{},error) {
	return ctx.container.Make(key)
}

func (ctx *Context) MustMake(key string) (interface{}) {
	return ctx.container.MustMake(key)
}

func (ctx *Context) MakeNew(key string,params []interface{}) (interface{},error) {
	return ctx.container.MakeNew(key,params)
}