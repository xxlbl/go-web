package gee

import (
	"log"
	"time"
)

// 自定义中间件
func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next() // 表示等待执行其他的中间件或用户的Handler
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
