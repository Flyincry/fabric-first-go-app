/**
  author: kevin
*/
package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
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
func (app *Application) ModifyShow(w http.ResponseWriter, r *http.Request) {
	// 根据证书编号与姓名查询信息
	jeweler := r.FormValue("jeweler")
	paperNumber := r.FormValue("paperNumber")
	action := r.FormValue("Action")
	result, err := app.Fabric.Querypaper(jeweler, paperNumber)

	var paper = service.InventoryFinancingPaper{}
	json.Unmarshal(result, &paper)

	data := &struct {
		Paper  service.InventoryFinancingPaper
		Msg    string
		Flag   bool
		Action string
	}{
		Paper:  paper,
		Msg:    "",
		Flag:   false,
		Action: action,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	showView(w, r, "setInfo.html", data)
}

func (app *Application) Modify(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	financingAmount, _ := strconv.Atoi(r.FormValue("financingAmount"))
	paper := service.InventoryFinancingPaper{
		Jeweler:            r.FormValue("jeweler"),
		PaperNumber:        r.FormValue("paperNumber"),
		FinancingAmount:    financingAmount,
		ApplyDateTime:      r.FormValue("applyDateTime"),
		ReviseDateTime:     r.FormValue("reviseDateTime"),
		AcceptDateTime:     r.FormValue("acceptDateTime"),
		ReadyDateTime:      r.FormValue("readyDateTime"),
		EvalDateTime:       r.FormValue("evalDateTime"),
		ReceiveDateTime:    r.FormValue("receiveDateTime"),
		EndDate:            r.FormValue("endDateTime"),
		PaidbackDateTime:   r.FormValue("paidBackDateTime"),
		RepurchaseDateTime: r.FormValue("RepurchaseDateTime"),
		Bank:               r.FormValue("bank"),
		Evaluator:          r.FormValue("evaluator"),
		Repurchaser:        r.FormValue("repurchaser"),
		Supervisor:         r.FormValue("supervisor"),
	}
	action := r.FormValue("Action")

	// 调用业务层, 反序列化
	app.Fabric.Action(paper, action)
	// result, err := app.Fabric.Action2(action, Jeweler, PaperNumber, FinancingAmount, ApplyDateTime, ReviseDateTime, AcceptDateTime, ReadyDateTime, EvalDateTime, ReceiveDateTime, EndDate, PaidbackDateTime, RepurchaseDateTime, Bank, Evaluator, Repurchaser, Supervisor)

	// // 封装响应数据
	// data := &struct {
	// 	Flag bool
	// 	Msg  string
	// 	Err  error
	// }{
	// 	Flag: true,
	// 	Msg:  result,
	// 	Err:  err,
	// }
	// if err != nil {
	// 	data.Msg = err.Error()
	// } else {
	// 	data.Msg = "操作成功，交易ID: " + transactionID
	// }

	// 响应客户端
	r.Form.Set("jeweler", paper.Jeweler)
	r.Form.Set("paperNumber", paper.PaperNumber)
	app.QueryInfo(w, r)
}

// 根据指定的 Key 查询信息
func (app *Application) QueryInfo(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	jeweler := r.FormValue("jeweler")
	paperNumber := r.FormValue("paperNumber")

	// 调用业务层, 反序列化
	msg, err := app.Fabric.Querypaper(jeweler, paperNumber)
	var paper = service.InventoryFinancingPaper{}
	json.Unmarshal(msg, &paper)

	// fmt.Println("查询信息成功：")
	// fmt.Println(paper)

	// 封装响应数据
	data := &struct {
		Paper service.InventoryFinancingPaper
		Msg   string
		Flag  bool
	}{
		Paper: paper,
		Msg:   "",
		Flag:  false,
	}
	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}
	// 响应客户端
	showView(w, r, "queryReq.html", data)
}

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
	showView(w, r, "other.html", data)
}
