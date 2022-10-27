package middleware

import (
	"context"
	"fmt"
	"goweb/framework/gin"
	"log"
	"time"
)

//超时的中间件
func TimeoutHandler(t time.Duration) gin.HandlerFunc {
	// 使用函数回调 
	return func (c *gin.Context) {
		finish := make(chan struct{},  1)
		panicChan := make(chan  interface{}, 1)
		//执行业务逻辑前预操作
		durationCtx, cancel := context.WithTimeout(c.Request.Context(),t)
		defer cancel()

		go func() {
			defer func (){
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			//执行具体的业务
			c.Next()

			finish <-struct{}{}
		}()
		// 执行业务逻辑后操作
		select{
		case p := <-panicChan:
			log.Println(p)
			c.ISetStatus(500).IJson("time out")
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.ISetStatus(500).IJson("Time Out")
		}
	}
}