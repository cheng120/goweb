package config

import (
	"bytes"
	"errors"
	"fmt"
	"goweb/framework"
	"goweb/framework/contract"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
)

type Config struct {
	c framework.Container
	folder string
	KeyDelim string // 路径的分隔符，默认为点
	lock sync.RWMutex
	envMaps map[string]string // 所有的环境变量
	confMaps map[string]interface{} // 配置文件结构，key为文件名
	confRaws map[string][]byte // 配置文件的原始信息
}

// 读取某个配置文件
func (c *Config) loadConfigFile(folder string, file string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	//  判断文件是否以yaml或者yml作为后缀
	s := strings.Split(file,".")
	if len(s) == 2 && (s[1] == "yml" || s[1] == "yaml" ){
		name := s[0]
		// 读取文件内容
		bf, err := os.ReadFile(filepath.Join(folder,file))
		if err != nil {
			return err
		}
		// 直接针对文本做环境变量的替换
		bf = replace(bf, c.envMaps)
		// 解析对应的文件
		conf := map[string]interface{}{}
		if err := yaml.Unmarshal(bf, &conf);err != nil {
			return err
		}
		c.confMaps[name] = conf
		c.confRaws[name] = bf

		// 读取app.path中的信息，更新app对应的folder
		if name == "app" && c.c.IsBind(contract.AppKey) {
			if p, ok := conf["path"]; ok {
				appService := c.c.MustMake(contract.AppKey).(contract.App)
				appService.LoadAppConfig(cast.ToStringMapString(p))
			}
		}
	}
	return nil
}


// 删除文件的操作
func (c *Config) removeConfigFile(folder string, file string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	s := strings.Split(file,".")
	// 只有yaml或者yml后缀才执行
	if len(s) == 2 && (s[1] == "yml" || s[1] == "yaml"){
		name := s[0]
		// 删除内存中对应的key
		delete(c.confRaws,name)
		delete(c.confMaps,name)
	}
	return nil
}


// NewConfig 初始化Config方法
func NewConfig(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	envFolder := params[1].(string)
	envMaps := params[2].(map[string]string)

	// 检查文件夹是否存在
	if _, err := os.Stat(envFolder);os.IsNotExist(err) {
		return nil, errors.New("folder " + envFolder + " not exist: " + err.Error())
	}
	// 实例化
	conf := &Config{
		c: container,
		envMaps: envMaps,
		folder: envFolder,
		confMaps: map[string]interface{}{},
		confRaws: map[string][]byte{},
		KeyDelim: ".",
		lock: sync.RWMutex{},
	}

	// 读取每个文件
	files, err := ioutil.ReadDir(envFolder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileName := file.Name()
		err := conf.loadConfigFile(envFolder, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	// 监控文件夹文件
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil,err
	}

	err = watch.Add(envFolder)
	if err != nil {
		return nil,err
	}

	go func() {
		defer func ()  {
			if err := recover();err != nil {
				fmt.Println(err)
			}
		}()

		for {
			select{
			case ev := <- watch.Events:
				{
					//判断事件发生的类型
					// Create 创建
					// Write 写入
					// Remove 删除
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator)) 	
					folder := path[:index]
					fileName := path[index+1:]
					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件 : ", ev.Name)
						conf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件 : ", ev.Name)
						conf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件 : ", ev.Name)
						conf.removeConfigFile(folder, fileName)
					}
				}
			case err := <- watch.Errors:
				{
					
					log.Println("error : ", err)
					return
				}
			}

		}
	}()
	return conf, nil
}


// replace 表示使用环境变量maps替换context中的env(xxx)的环境变量
func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}
	// 直接使用ReplaceAll替换。这个性能可能不是最优，但是配置文件加载，频率是比较低的，可以接受
	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content,[]byte(reKey),[]byte(val))
	}
	return content
}


// 查找某个路径的配置项
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}
	// 判断是否有下个路径
	next, ok := source[path[0]]
	if ok {
		if len(path) == 1 {
			return next
		}
		// 判断下一个路径的类型
		switch next.(type) {
		case map[interface{}]interface{}:
			// 如果是interface的map，使用cast进行下value转换
			return searchMap(cast.ToStringMap(next),path[1:])
		case map[string]interface{}:
			// 如果是map[string]，直接循环调用
			return searchMap(next.(map[string]interface{}),path[1:])
		default:
			return nil
		}
	}
	return nil
}

// 通过path来获取某个配置项
func (conf *Config) find(key string) interface{} {
	conf.lock.RLock()
	defer conf.lock.RUnlock()
	return searchMap(conf.confMaps, strings.Split(key, conf.KeyDelim))
}

// IsExist check setting is exist
func (conf *Config) IsExist(key string) bool {
	return conf.find(key) != nil
}

// Get 获取某个配置项
func (conf *Config) Get(key string) interface{} {
	return conf.find(key)
}

// GetBool 获取bool类型配置
func (conf *Config) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

// GetInt 获取int类型配置
func (conf *Config) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

// GetFloat64 get float64
func (conf *Config) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}

// GetTime get time type
func (conf *Config) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}

// GetString get string typen
func (conf *Config) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

// GetIntSlice get int slice type
func (conf *Config) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}

// GetStringSlice get string slice type
func (conf *Config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

// GetStringMap get map which key is string, value is interface
func (conf *Config) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(conf.find(key))
}

// GetStringMapString get map which key is string, value is string
func (conf *Config) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

// GetStringMapStringSlice get map which key is string, value is string slice
func (conf *Config) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

// Load a config to a struct, val should be an pointer
func (conf *Config) Load(key string, val interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(conf.find(key))
}