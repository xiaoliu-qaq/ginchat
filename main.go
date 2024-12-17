// @Author TrandLiu
// @Date 2024/12/12 17:22:00
// @Desc
package main

import (
	"newginchat/router"
	"newginchat/utils"
)

func main() {
	//连接数据库
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	r := router.Router()
	r.Run(":8081")
}
