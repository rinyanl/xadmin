package main

import (
	"xadmin/app"
)

func main() {
	// go socket.Test()
	app.RunApp()
	// 同时添加多个用户、防止频繁重启 xary 测试
	// for i := 0; i < 200; i++ {
	// users.XaryAddUser(fmt.Sprint(i) + "aa@qq.com")
	// }
	// time.Sleep(10 * time.Second)
}
