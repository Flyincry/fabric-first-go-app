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

	http.HandleFunc("/", app.IndexView)
	http.HandleFunc("/index.html", app.IndexView)
	http.HandleFunc("/setInfo.html", app.SetInfoView)
	http.HandleFunc("/setReq", app.Apply)
	http.HandleFunc("/queryReq", app.QueryInfo)
	http.HandleFunc("/query", app.QueryInfo)

	fmt.Println("启动Web服务, 监听端口号: 9001")

	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}

}
