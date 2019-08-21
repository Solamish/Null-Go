# Null-Go FrameWork
> 使用过Golang的coder应该都知道gin这款web框架，其简洁的设计以及高效的性能和良好的生产力，让使用过的人都爱不释手。于是，我抱着学习的态度，模仿着gin的思路实现了一款简单的web框架。
## 快速开始
```go
package main

import (
	"nullgo"
)

func main() {
	router := nullgo.Default()
	router.GET("/hello", func(c *nullgo.Context) {
		c.String("hello,world")
	})
	router.Run(":8080")
}
```

## 路由
* 使用GET,POST,PUT,DELETE
```go
package main

import (
	"nullgo"
)

func main() {
        // Initialize a null-go router
	router := nullgo.Default()

	router.GET("/someGet", func(context *nullgo.Context) {})
	router.POST("/somePost", func(context *nullgo.Context) {})
	router.PUT("/somePut", func(context *nullgo.Context) {})
	router.DELETE("/someDelete", func(context *nullgo.Context) {})
	
	//The example below serves on port :8080
	//Also, you can serve it on any port as you want
	router.Run(":8080")
}
```

* 获取GET参数
```go
func main() {
	router := nullgo.Default()

	router.GET("/user", func(c *nullgo.Context) {
		name := c.Query("name")
		age := c.Query("age")
		c.String("name:" + name)
		c.String("age:"+age)
	})
}
```

* 获取POST参数
```go
func main() {
  	router := nullgo.Default()
  
  	router.POST("/hello", func(c *nullgo.Context) {
  		message := c.PostV("message")
  		c.String("message:" + message)
  	})
  }
```

* 获取路径中的参数
```go
func main() {
	router := nullgo.Default()

	router.GET("/hello/:name", func(c *nullgo.Context) {
		name := c.Param("name")
		c.String("hello" + name)
	})
}
```

* 使用正则
```go
func main() {
	router := nullgo.Default()
	
	roouter.GET("/user/:id([1-9]+)", func(c *nullgo.Context) {
		id := c.Param("id)
		c.String("id:",id)
	}
}
```


## Log 提供的一些方法
```go
func Trace(format string, v ...interface{}) {
	if level <= LevelTrace {
		NullLogger.Printf("[BULL-trace] "+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if level <= LevelTrace {
		NullLogger.Printf("[BULL-error] "+format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if level <= LevelTrace {
		NullLogger.Printf("[BULL-warn] "+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if level <= LevelTrace {
		NullLogger.Printf("[BULL-info] "+format, v...)
	}
}
func Debug(format string, v ...interface{}) {
	if level <= LevelTrace {
		NullLogger.Printf("[BULL-debug] "+format, v...)
	}
}
```


## 读取配置文件
* config 
```sh
#database config
username = root
password  = root
DBName = redrock
```

* load
```go
func GetInfo() {
	cfg, err := LoadConfig("test.conf")
	if err != nil {
		Error("get info error")
		return
	}
	name := cfg.String("username")
	pwd := cfg.String("password")
	DBName := cfg.String("DBName")

	//fmt.Println("name:", name)
	//fmt.Println("password:", pwd)
	//fmt.Println("DataBaseName:", DBName)
}
```

## Release History 版本历史
* v 0.1.0


##  TODO 
- [ ] 使用树实现路由
- [ ] WebSocket完善
- [ ] Middleware实现
- [ ] 复杂路由支持(路由组..)
- [ ] 安全机制
- [ ] ...
