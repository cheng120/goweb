package framework

import (
	"errors"
	"fmt"
	"sync"
)

// Container 是一个服务容器，提供服务和获取服务的功能
type Container interface {
	//Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，返回error
	Bind(provider ServiceProvider) error
	//IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	//Make 根据关键字提供一个服务
	Make(key string) (interface{}, error)

	//MustMake 根据关键字获取一个服务，如果这个关键字未绑定服务提供者，返回panic
	//所以在使用这个接口的时候保证已经绑定服务
	MustMake(key string) interface{}

	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的params参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// MainContainer 服务容器的具体实现
type ServiceContainer struct {
	Container
	//providers 存储注册的服务提供者，key为字符串凭证
	providers map[string]ServiceProvider
	// instance 存储具体的实例 ,key为凭证
	instances map[string]interface{}
	// lock 用于锁住对容器的变更操作 读写锁
	lock sync.RWMutex
}

//创建一个服务容器
func NewContainer() *ServiceContainer {
	return &ServiceContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock: sync.RWMutex{},
	}
}

//PrintProviders 输出服务容器中注册的关键字
func (sc *ServiceContainer) PrintProviders() []string {
	ret := []string{}
	for _,provider := range sc.providers { 
		name := provider.Name()
		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	return ret
}

// Bind 将服务容器和关键字做了绑定
func (sc *ServiceContainer) Bind(prodiver ServiceProvider) error {
	sc.lock.Lock()
	key := prodiver.Name()
	sc.providers[key] = prodiver

	sc.lock.Unlock()
	// if provider is not defer
	if !prodiver.IsDefer()  {
		if err := prodiver.Boot(sc);err != nil {
			return err
		}
		//实例化
		params := prodiver.Params(sc)
		method := prodiver.Register(sc)
		instance, err := method(params...)
		if err != nil {
			fmt.Println("bind service provider ", key, " error: ", err)
			return errors.New(err.Error())
		}

		sc.instances[key] = instance
	}

	return nil
}

func (sc *ServiceContainer) IsBind(key string) bool {
	return sc.findServiceProvider(key) != nil
}
func (sc *ServiceContainer) findServiceProvider(key string) ServiceProvider {
	sc.lock.RLock()
	defer sc.lock.RUnlock()	
	if sp,ok := sc.providers[key];ok {
		return sp
	}
	return nil
}

func (sc *ServiceContainer) Make(key string) (interface{}, error) {
	return sc.make(key,nil,false)
}

func (sc *ServiceContainer) MustMake(key string) (interface{})  {
	sp,err := sc.make(key,nil,true)
	if err != nil {
		panic(err)
	}
	

	return sp
}

func (sc *ServiceContainer) MakeNew(key string, params[]interface{}) (interface{}, error) {
	return sc.make(key, params, true)
}

func (sc *ServiceContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	
	//force a new 
	if err := sp.Boot(sc);err != nil {
		return nil,err
	}
	if params == nil {
		params = sp.Params(sc)
	}
	method := sp.Register(sc)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins,nil
}


func (sc *ServiceContainer) make(key string,params []interface{}, forceNew bool) (interface{}, error) {
	sc.lock.RLock()
	defer sc.lock.RUnlock()
	// 查询是否已经注册了这个服务提供者，如果没有注册，则返回错误
	sp := sc.findServiceProvider(key)

	if sp == nil {
		return nil,errors.New("contract " + key + " have not register")
	}

	
	// 不需要强制重新实例化，如果容器中已经实例化了，那么就直接使用容器中的实例
	
	if ins, ok := sc.instances[key];ok {
		return ins,nil
	}

	if forceNew {
		return sc.newInstance(sp,params)
	}
	// 容器中还未实例化，则进行一次实例化
	insn,err := sc.newInstance(sp,params)
	if err != nil {
		return nil, err
	}
	sc.instances[key] = insn
	return insn,nil
}


// NameList 列出容器中所有服务提供者的字符串凭证
func (c *ServiceContainer) NameList() []string {
	ret := []string{}
	for _, provider := range c.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}