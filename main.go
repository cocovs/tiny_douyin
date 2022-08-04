package main

import (
	"github.com/cocovs/tiny-douyin/setting"
	"github.com/gin-gonic/gin"
)

func main() {
	//加载配置文件
	setting.Setini()

	//主函数最后关闭数据库，否则可能会意外关闭
	//defer mysql.Close()
	//默认引擎
	r := gin.Default()
	//路由入口
	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
