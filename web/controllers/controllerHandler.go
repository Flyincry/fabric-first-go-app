/**
  author: kevin
*/
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/db/model"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/sdkenv"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/service"
)

var store = sessions.NewCookieStore([]byte("mysession"))

type Application struct {
	Fabric *service.ServiceHandler
	SDK    *fabsdk.FabricSDK
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
		var err error
		go func() {
			err = app.Fabric.CreateOrg(u1.OrgID)
		}()
		if err != nil {
			data := map[string]interface{}{
				"err": "创建组织失败:" + err.Error(),
			}
			showView(w, r, "Regist.html", data)
		} else {
			fmt.Println("已添加如下用户", u1)
			showView(w, r, "login.html", nil)
		}
	}

	// 这需要添加处理注册错误的程序
}

func (app *Application) Account(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "login.html", nil)
}

func (app *Application) Login_cn(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "login_cn.html", nil)
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
			session.Values["orgid"] = u.OrgID
			session.Values["channel"] = request.Form.Get("channel")
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
	http.Redirect(response, request, "/home", http.StatusSeeOther)
}

// 根据指定的 key 设置/修改 value 信息
func (app *Application) ModifyShow(w http.ResponseWriter, r *http.Request) {
	// 根据证书编号与姓名查询信息
	jeweler := r.FormValue("jeweler")
	paperNumber := r.FormValue("paperNumber")
	action := r.FormValue("Action")
	paper := service.InventoryFinancingPaper{
		Jeweler:     jeweler,
		PaperNumber: paperNumber,
	}
	msg, err := app.Fabric.Action(paper, "QueryPaper")
	//fmt.Println(msg)
	json.Unmarshal(msg, &paper)

	session, _ := store.Get(r, "mysession")
	username := session.Values["username"].(string)
	role := session.Values["role"].(string)

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
	var paper service.InventoryFinancingPaper
	var decoder = schema.NewDecoder()
	decoder.SetAliasTag("json")
	err := r.ParseForm()
	err = decoder.Decode(&paper, r.Form)
	fmt.Println(paper)
	fmt.Println("paper state: ", paper.State, "\n")

	// 执行链码
	action := r.FormValue("Action")
	fmt.Println(action)

	// 创建服务句柄, 调用业务层, 反序列化
	session, _ := store.Get(r, "mysession")
	org := sdkenv.OrgInfo{
		OrgName: "Org" + session.Values["orgid"].(string),
		OrgUser: "User1",
	}
	//channel, _ := session.Values["channel"].(string)
	//fmt.Println("Modify channel:" + channel + "\n")
	fmt.Println("Modify org:" + org.OrgName + org.OrgUser + "\n")
	//serviceHandler, err := service.InitService("simplecc", channel, &org, app.SDK)
	//res, err := serviceHandler.Action(paper, action)
	_, err = app.Fabric.Action(paper, action)
	//fmt.Println(res)
	fmt.Println(err)

	// 响应客户端
	r.Form.Set("jeweler", paper.Jeweler)
	r.Form.Set("paperNumber", paper.PaperNumber)
	app.QueryInfo1(w, r)
}

func (app *Application) QueryShow(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"]
	role := session.Values["role"]
	data := map[string]interface{}{
		"username": username,
		"role":     role,
	}
	showView(w, r, "queryPage.html", data)
}

// 根据指定的 Key 查询信息test
func (app *Application) QueryInfo1(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	jeweler := r.FormValue("jeweler")
	paperNumber := r.FormValue("paperNumber")
	paper := service.InventoryFinancingPaper{
		Jeweler:     jeweler,
		PaperNumber: paperNumber,
	}
	action := "QueryPaper"

	// 获取用户信息
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"].(string)
	role := session.Values["role"].(string)
	// fmt.Println("username: ", username)
	// fmt.Println("role: ", role)

	// 调用业务层, 反序列化
	//org := sdkenv.OrgInfo{
	//OrgName: "Org" + session.Values["orgid"].(string),
	//OrgUser: "User1",
	//}
	//channel, _ := session.Values["channel"].(string)
	//serviceHandler, err := service.InitService("simplecc", channel, &org, app.SDK)
	//msg, err := serviceHandler.Action(paper, action)
	msg, err := app.Fabric.Action(paper, action)
	//fmt.Println(msg)
	json.Unmarshal(msg, &paper)
	fmt.Println(paper)
	fmt.Println("paper state: ", paper.State, "\n")

	// 封装响应数据
	data := &struct {
		Paper    service.InventoryFinancingPaper
		Msg      string
		Flag     bool
		Username string
		Role     string
		Action   string
	}{
		Paper:    paper,
		Msg:      "",
		Flag:     false,
		Username: username,
		Role:     role,
		Action:   action,
	}
	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}
	// 响应客户端
	showView(w, r, "queryResult.html", data)
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
	go app.Fabric.CreateChan(ChannelName)

	// 响应客户端
	// r.Form.Set("OrgName", "1")
	// r.Form.Set("Port", "7051")
	// app.QueryChannel(w, r)
	app.CreateChannelShow(w, r)
}

// 根据指定的 Key 查询某个组织所在的通道
func (app *Application) QueryChannel(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	OrgName := r.FormValue("OrgName")

	// 获取用户信息
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"].(string)
	role := session.Values["role"].(string)

	msg := ""
	// 调用业务层, 反序列化
	if OrgName != "" {
		msg, _ = app.Fabric.QueryChan(OrgName)
	}

	// 封装响应数据
	data := &struct {
		Msg      string
		Username string
		Role     string
	}{
		Msg:      msg,
		Username: username,
		Role:     role,
	}

	// if err != nil {
	// 	data.Msg = err.Error()
	// }

	// 响应客户端
	showView(w, r, "QueryChannel.html", data)
}

// 展示加入channel的界面
func (app *Application) JoinChannelShow(w http.ResponseWriter, r *http.Request) {
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

	showView(w, r, "JoinChannel.html", data)
}

// 让某个组织加入某个通道
func (app *Application) JoinChannel(w http.ResponseWriter, r *http.Request) {
	// 获取提交数据
	OrgName := r.FormValue("OrgName")
	ChannelName := r.FormValue("ChannelName")

	// 获取用户信息
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"].(string)
	role := session.Values["role"].(string)

	msg := ""
	// 调用业务层, 反序列化
	if OrgName != "" {
		msg, _ = app.Fabric.JoinChan(OrgName, ChannelName)
	}

	// 封装响应数据
	data := &struct {
		Msg      string
		Username string
		Role     string
	}{
		Msg:      msg,
		Username: username,
		Role:     role,
	}

	// if err != nil {
	// 	data.Msg = err.Error()
	// }

	// 响应客户端
	showView(w, r, "JoinChannel.html", data)
}
