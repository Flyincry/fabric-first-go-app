/**
  author: kevin
*/
package controllers

import (
	"net/http"
	"time"

	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/service"
)

type Application struct {
	Fabric *service.ServiceHandler
}

func (app *Application) IndexView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "index.html", nil)
}

func (app *Application) SetInfoView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "setInfo.html", nil)
}

// 根据指定的 key 设置/修改 value 信息
func (app *Application) Apply(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	jeweler := r.FormValue("jeweler")
	paperNumber := r.FormValue("paperNumber")
	financialAmount := r.FormValue("financialAmount")
	applyDateTime := time.Now().String()

	// 调用业务层, 反序列化
	transactionID, err := app.Fabric.Apply(paperNumber, jeweler, applyDateTime, financialAmount)

	// 封装响应数据
	data := &struct {
		Flag bool
		Msg  string
	}{
		Flag: true,
		Msg:  "",
	}
	if err != nil {
		data.Msg = err.Error()
	} else {
		data.Msg = "操作成功，交易ID: " + transactionID
	}

	// 响应客户端
	showView(w, r, "setInfo.html", data)
}

// 根据指定的 Key 查询信息
func (app *Application) QueryInfo(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	jeweler := r.FormValue("jeweler")
	paperNumber := r.FormValue("paperNumber")

	// 调用业务层, 反序列化
	msg, err := app.Fabric.Querypaper(jeweler, paperNumber)

	// 封装响应数据
	data := &struct {
		Msg string
	}{
		Msg: "",
	}
	if err != nil {
		data.Msg = "没有查询到对应的信息"
	} else {
		data.Msg = "查询成功: " + msg
	}
	// 响应客户端
	showView(w, r, "queryReq.html", data)
}
