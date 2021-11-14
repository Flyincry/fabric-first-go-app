/**
  author: kevin
*/

// 路由文件
package web

import (
	"fmt"
	"net/http"

	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/web/controllers"
)

func WebStart(app *controllers.Application) {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// http.HandleFunc("/index", app.IndexView)
	http.HandleFunc("/index.html", app.IndexView)
	http.HandleFunc("/query", app.QueryInfo)
	http.HandleFunc("/modifyPage", app.ModifyShow)
	http.HandleFunc("/modify", app.Modify)
	//http.HandleFunc("/Apply", app.Apply)

	http.HandleFunc("/", app.Account)
	http.HandleFunc("/home", app.Account)
	http.HandleFunc("/home/index", app.Account)
	http.HandleFunc("/home/login", app.Login)
	http.HandleFunc("/home/welcome", app.Welcome)
	http.HandleFunc("/home/logout", app.Logout)
	// http.HandleFunc("/home/logout", accountcontroller.Logout)

	http.HandleFunc("/QueryChannel", app.QueryChannel)
	http.HandleFunc("/CreateChannelShow", app.CreateChannelShow)
	http.HandleFunc("/CreateChannel", app.CreateChannel)

	http.HandleFunc("/RegistPage", app.RegistPage)
	http.HandleFunc("/Regist", app.Regist)

	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}

}
