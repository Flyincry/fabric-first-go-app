/**
  author: kevin
*/
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/db/model"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/service"
)

var store = sessions.NewCookieStore([]byte("mysession"))

type Application struct {
	Fabric *service.ServiceHandler
}

func (app *Application) RegistPage(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "Regist.html", nil)
}

func (app *Application) Regist(w http.ResponseWriter, r *http.Request) {
	u := model.User{
		Name:         r.FormValue("username"),
		Password:     r.FormValue("password"),
		Role:         r.FormValue("role"),
		Organization: r.FormValue("organization"),
	}
	// 检查用户是否存在
	user, _ := u.GetUserByName()
	if user.Name != "" {
		//用户名已存在
		data := map[string]interface{}{
			"err": "创建用户失败：用户名已存在",
		}
		showView(w, r, "Regist.html", data)
	} else {
		_ = u.AddUser()
		u1, _ := u.GetUserByName()
		err := app.Fabric.CreateOrg(u1.OrgID)
		if err != nil {
			data := map[string]interface{}{
				"err": "创建组织失败:" + err.Error(),
			}
			showView(w, r, "Regist.html", data)
		} else {
			fmt.Println("已添加如下用户", u1)
			showView(w, r, "Regist.html", nil)
		}
	}

	// 这需要添加处理注册错误的程序
}

func (app *Application) Account(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "login.html", nil)
}

func (app *Application) Login(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	username := request.Form.Get("username")
	password := request.Form.Get("password")
	u := &model.User{
		Name: username,
	}
	u, _ = u.GetUserByName()

	if u != nil {
		//存在
		if password == u.Password {
			session, _ := store.Get(request, "mysession")
			session.Values["username"] = username
			session.Values["role"] = u.Role
			session.Save(request, response)
			http.Redirect(response, request, "/home/welcome", http.StatusSeeOther)
		} else {
			data := map[string]interface{}{
				"err": "请输入与用户名相匹配的密码",
			}
			showView(response, request, "login.html", data)
		}
	} else {
		data := map[string]interface{}{
			"err": "用户名不存在，请输入合法用户名",
		}
		showView(response, request, "login.html", data)
	}
}

func (app *Application) Welcome(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "mysession")
	username := session.Values["username"]
	role := session.Values["role"]
	fmt.Println("username: ", username)
	fmt.Println("role: ", role)
	data := map[string]interface{}{
		"username": username,
		"role":     role,
	}
	showView(response, request, "welcome.html", data)
}

func (app *Application) Logout(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "mysession")
	session.Options.MaxAge = -1
	session.Save(request, response)
	http.Redirect(response, request, "/home/index", http.StatusSeeOther)
}

func (app *Application) IndexView(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"]
	role := session.Values["role"]
	data := map[string]interface{}{
		"username": username,
		"role":     role,
	}
	showView(w, r, "index.html", data)
}

// func (app *Application) SetInfoView(w http.ResponseWriter, r *http.Request) {
// 	showView(w, r, "setInfo.html", nil)
// }

// 根据指定的 key 设置/修改 value 信息
func (app *Application) ModifyShow(w http.ResponseWriter, r *http.Request) {
	// 根据证书编号与姓名查询信息
	jeweler := r.FormValue("jeweler")
	paperNumber := r.FormValue("paperNumber")
	action := r.FormValue("Action")
	result, err := app.Fabric.Querypaper(jeweler, paperNumber)

	session, _ := store.Get(r, "mysession")
	username := session.Values["username"].(string)
	role := session.Values["role"].(string)

	var paper = service.InventoryFinancingPaper{}
	json.Unmarshal(result, &paper)

	data := &struct {
		Paper    service.InventoryFinancingPaper
		Msg      string
		Flag     bool
		Action   string
		Username string
		Role     string
	}{
		Paper:    paper,
		Msg:      "",
		Flag:     false,
		Action:   action,
		Username: username,
		Role:     role,
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
		RepurchaseDateTime: r.FormValue("repurchaseDateTime"),
		Bank:               r.FormValue("bank"),
		Evaluator:          r.FormValue("evaluator"),
		Repurchaser:        r.FormValue("repurchaser"),
		Supervisor:         r.FormValue("supervisor"),
	}
	action := r.FormValue("Action")

	// 调用业务层, 反序列化
	app.Fabric.Action(paper, action)

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

	// 获取用户信息
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"].(string)
	role := session.Values["role"].(string)
	// fmt.Println("username: ", username)
	// fmt.Println("role: ", role)

	// 调用业务层, 反序列化
	msg, err := app.Fabric.Querypaper(jeweler, paperNumber)
	var paper = service.InventoryFinancingPaper{}
	json.Unmarshal(msg, &paper)

	// 封装响应数据
	data := &struct {
		Paper    service.InventoryFinancingPaper
		Msg      string
		Flag     bool
		Username string
		Role     string
	}{
		Paper:    paper,
		Msg:      "",
		Flag:     false,
		Username: username,
		Role:     role,
	}
	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}
	// 响应客户端
	showView(w, r, "queryReq.html", data)
}

// Channel
func (app *Application) CreateChannelShow(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"].(string)
	role := session.Values["role"].(string)

	data := &struct {
		Username string
		Role     string
	}{
		Username: username,
		Role:     role,
	}

	showView(w, r, "CreateChannelShow.html", data)
}

func (app *Application) CreateChannel(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	ChannelName := r.FormValue("ChannelName")

	// 调用业务层, 反序列化
	app.Fabric.CreateChan(ChannelName)

	// 响应客户端
	r.Form.Set("OrgName", "1")
	r.Form.Set("Port", "7051")
	app.QueryChannel(w, r)
}

// 根据指定的 Key 查询信息
func (app *Application) QueryChannel(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	OrgName := r.FormValue("OrgName")
	Port := r.FormValue("Port")

	// 获取用户信息
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"].(string)
	role := session.Values["role"].(string)

	msg, err := "", ""
	// 调用业务层, 反序列化
	if OrgName != "" {
		msg = app.Fabric.QueryChan(OrgName, Port)
	}

	// 封装响应数据
	data := &struct {
		Msg      string
		Err      string
		Username string
		Role     string
	}{
		Msg:      msg,
		Err:      err,
		Username: username,
		Role:     role,
	}

	// 响应客户端
	showView(w, r, "QueryChannel.html", data)
}

// func (app *Application) Apply(w http.ResponseWriter, r *http.Request) {
// 	// 获取提交数据
// 	jeweler := r.FormValue("jeweler")
// 	paperNumber := r.FormValue("paperNumber")
// 	financialAmount := r.FormValue("financialAmount")
// 	applyDateTime := time.Now().String()

// 	// 调用业务层, 反序列化
// 	transactionID, err := app.Fabric.Apply(paperNumber, jeweler, applyDateTime, financialAmount)

// 	// 封装响应数据
// 	data := &struct {
// 		Flag bool
// 		Msg  string
// 	}{
// 		Flag: true,
// 		Msg:  "",
// 	}
// 	if err != nil {
// 		data.Msg = err.Error()
// 	} else {
// 		data.Msg = "操作成功，交易ID: " + transactionID
// 	}

// 	// 响应客户端
// 	showView(w, r, "other.html", data)
// }
