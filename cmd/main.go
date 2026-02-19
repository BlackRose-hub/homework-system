package main

//程序的入口，负责启动HTTP服务器
import (
	"homework-system/configs" //导入配置包
	"homework-system/router"  //导入路由包
	"log"                     //日志记录
)

func main() {
	configs.Init()                         //这里用于加载配置、连接数据库是，加载配置，连接Mysql数据库，设置时区，初始化日志系统等
	r := router.SetupRouter()              //设置路由，注册所有API接口
	if err := r.Run(":8080"); err != nil { //启动服务器，监听8080接口HTTP端口）
		log.Fatal("服务器启动失败：", err) //如果启动失败，记录错误并退出程序
	}
}
