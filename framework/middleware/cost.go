package middleware

import (
	"goweb/framework/gin"
	"log"
	"time"
)

//recover 机制，将协程中的异常捕获
func Cost() gin.HandlerFunc {
	//使用函数回调
	return func(c *gin.Context)  {
		// 记录开始时间
		start := time.Now()

		// 使用next执行具体逻辑
		c.Next()
		//记录结束时间
		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri: %v, cost: %v", c.Request.RequestURI, cost.Seconds())
	}
}