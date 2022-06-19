/**
  author: kevin
*/
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "net"
	"log"
	// "time"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"errors"
	"strconv"
	"database/sql"
	_ "github.com/lib/pq"

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

// 解析公匙
func parsingRsaPublicKey(file string) (*rsa.PublicKey, error) {
	// 读取公匙文件
	pubByte,err := ioutil.ReadFile(file)
	if err != nil {
	  return nil,err
	}
	// pem解码
	b,_ := pem.Decode(pubByte)
	if b == nil {
	  return nil,errors.New("error public key")
	}
	// der解码，最终返回一个公匙对象
	pubKey,err := x509.ParsePKCS1PublicKey(b.Bytes)
	if err != nil {
	  return nil,err
	}
	return pubKey,nil
  }
  
//  rsa公匙加密
func rsaPublicKeyEncrypt(src []byte, publickey *rsa.PublicKey) ([]byte, error) {
	// 使用公匙加密数据，需要一个随机数生成器和公匙和需要加密的数据
	data,err := rsa.EncryptPKCS1v15(rand.Reader, publickey,src)
	if err != nil {
	  return nil,err
	}
	return data,nil
  }
  // rsa私匙解密
  
func rsaPrivateKeyDecrypt(src []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	// 使用私匙解密数据，需要一个随机数生成器和私匙和需要解密的数据
	data,err := rsa.DecryptPKCS1v15(rand.Reader, privateKey,src)
	if err != nil {
	  return nil,err
	}
	return data,nil
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

	// 取消从这里执行Action
	_, err = app.Fabric.Action(paper, action)
	//// fmt.Println(res)
	fmt.Println(err)

	// // ---------------------------------
	// // start a client
	// // 根据不同的action来决定接口
	// var service string
	// if action == "OfferProductInfo"{
	// 	service = "0.0.0.0:8900"
	// }else if action == "OfferLisenceInfo"{
	// 	service = "0.0.0.0:8901"
	// }else if action == "Evaluate"{
	// 	service = "0.0.0.0:8902"
	// }else if action == "PutInStorage"{
	// 	service = "0.0.0.0:8903"
	// }else if action == "ReadyRepo"{
	// 	service = "0.0.0.0:8904"
	// }else{
	// 	service = "0.0.0.0:8888"
	// }
	// // service := "0.0.0.0:8888"
    // fmt.Println(service)
    // tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    // if err != nil {
    //     // log.Fatal(err)
	// 	log.Print("The server is down, it should be", service)
    // }
    // conn, err := net.DialTCP("tcp4", nil, tcpAddr)
    // if err != nil {
    //     log.Fatal(err)
    // }

	// // fmt.Println(paper)
    // b, err := json.Marshal(paper)
    // if err != nil{
    //     fmt.Println("error:", err)
    // }
	// data := []byte(b)

	// // rsa 加密
	// pubKey,_ := parsingRsaPublicKey("pub.key") // 解密公匙
	// encryData,_ := rsaPublicKeyEncrypt(data,pubKey) // 加密数据
	// fmt.Printf("%x\n",encryData)
    
	// n, err := conn.Write(encryData)
    // if err != nil {
    //     log.Fatal(err)
    // }
	// fmt.Println(data)

	// time.Sleep(time.Duration(2)*time.Second)

	// Action := []byte(action)
    // nn, err := conn.Write(Action)
    // if err != nil {
    //     log.Fatal(err)
    // }
    // fmt.Println(n, nn)

    // // log.Fatal(n)
	// fmt.Println(Action)
	
	// // -------------------------------

	// 响应客户端
	r.Form.Set("jeweler", paper.Jeweler)
	r.Form.Set("paperNumber", paper.PaperNumber)
	// time.Sleep(time.Duration(3)*time.Second)
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

// postgres
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "gyy"
	dbname   = "postgres"
)

func connectDB() *sql.DB{
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Postgres successfully connected!")
	return db
}

func query_marketprice(db *sql.DB, Pledgetype string) string {
	var price string
	// select jewelry_price from marketprice where jewelry_name='K金'
	err := db.QueryRow("select jewelry_price from marketprice where jewelry_name=$1", Pledgetype).Scan(&price)

	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			fmt.Println(err)
			fmt.Println("请输入正确的Type")
			log.Fatal(err)
		}
	}

	// fmt.Println("market price is :", price)
	return price
}

func query_single_addr(db *sql.DB, jeweler_name string) string{
	var addr string
    err := db.QueryRow(" select jeweler_addr from jeweler where jeweler_name=$1",jeweler_name).Scan(&addr)
    if err != nil {
        panic(err)
    }
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			log.Fatal(err)
		}
	}

	// fmt.Println("market price is :", price)
	return addr
}

func query_his_addr(db *sql.DB) ([]string, []string){
	var name, addr string
    rows,err:=db.Query("select jeweler_name, jeweler_addr from jeweler")
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    var his_name []string
	var his_addr []string
    for rows.Next(){
        rows.Scan(&name, &addr)
        his_name = append(his_name, name)
		his_addr = append(his_addr, addr)
    }
    return his_name, his_addr
}

func GetMarketprice(Pledgetype string) float64{
	db:=connectDB()
	price := query_marketprice(db, Pledgetype)
	market_price, _ := strconv.ParseFloat(price,64)
	return market_price
}

// levenshtein算法实现字符串相似度匹配
func levenshtein(str1, str2 string) float64 {
	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= s1len; y++ {
			oldkey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = minimum(column[y]+1, column[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	var res float64
	res = float64(column[s1len]) / float64(len(str2))
	return 1 - res
}

func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}

func CheckFakeaddress(jeweler_name string, jeweler_addr string) bool{
	db:=connectDB()
	his_name, his_addr := query_his_addr(db)
	for i := 0; i < len(his_name); i++{
		if jeweler_name == his_name[i]{
			single_add := query_single_addr(db, jeweler_name)
			if single_add != jeweler_addr{
				return true
			}
		}else if levenshtein(his_addr[i], jeweler_addr) > 0.8 {
			return true
		}
	}
	return false
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
	fmt.Println("This is the query result of ", action)
	fmt.Println(paper)
	fmt.Println("paper state: ", paper.State, "\n")
	
	//--------------------------------------------------
	// 添加前端的fraud detection
	// fraud 变量来记录是否检测出fraud
	fraud := "fine"

	PledgeApraisedValue, err := strconv.ParseFloat(paper.PledgeApraisedValue,64)
	if err != nil {
		fmt.Println("PledgeApraisedValue 字符串转换失败")
	}
	PledgeAmount, err := strconv.ParseFloat(paper.PledgeAmount,64)
	if err != nil {
		fmt.Println("PledgeAmount 字符串转换失败")
	}
	MarketPrice := GetMarketprice(paper.PledgeType)
	var PledgeValue float64 = MarketPrice * PledgeAmount
	fmt.Println(PledgeApraisedValue, PledgeAmount, MarketPrice, PledgeValue)

	// 银行 receive
	if role == "Bank" && paper.Bank == ""{
		// invoicing
		PledgeApraisedValue, err := strconv.ParseFloat(paper.PledgeApraisedValue,64)
		if err != nil {
			fmt.Println("PledgeApraisedValue 字符串转换失败")
		}
		PledgeAmount, err := strconv.ParseFloat(paper.PledgeAmount,64)
		if err != nil {
			fmt.Println("PledgeAmount 字符串转换失败")
		}
		MarketPrice := GetMarketprice(paper.PledgeType)
		var PledgeValue float64 = MarketPrice * PledgeAmount
		if PledgeApraisedValue > 1.001 * PledgeValue {
			fraud = "Overinvoicing"
		} else if PledgeApraisedValue < 0.999 * PledgeValue{
			fraud = "Underinvoicing"
		}
		// fake address
		// fmt.Println(CheckFakeaddress(paper.Jeweler, paper.JewelerAddr))
		if CheckFakeaddress(paper.Jeweler, paper.JewelerAddr) == true{
			fraud = "Fakeaddress"
		}
	}

	// 银行 accept
	if role == "Bank" && paper.Supervisor != ""{
		// phantom shipping
		PledgeAmount, err := strconv.Atoi(paper.PledgeAmount)
		if err != nil{
			fmt.Println("PledgeAmount 字符串转换失败")
		}
		ProductAmount, err := strconv.Atoi(paper.ProductAmount)
		if err != nil{
			fmt.Println("ProductAmount 字符串转换失败")
		}
		StorageAmount, err := strconv.Atoi(paper.StorageAmount)
		if err != nil{
			fmt.Println("StorageAmount 字符串转换失败")
		}
		if PledgeAmount != ProductAmount || ProductAmount != StorageAmount{
			fraud = "Phantomshipping"
		}

		// check product quality
		PledgeType := paper.PledgeType
		EvalType := paper.EvalType
		EvalPrice, err := strconv.Atoi(paper.EvalPrice)
		PledgeApraisedValue, err := strconv.Atoi(paper.PledgeApraisedValue)
		if err != nil {
			fmt.Println("PledgeApraisedValue 字符串转换失败")
		}
		if err != nil{
			fmt.Println("EvalPrice 字符串转换失败")
		}
		if PledgeType != EvalType || PledgeApraisedValue <= int(0.999 * float64(EvalPrice)){
			fraud = "Checkproductquality"
		}
	}
	//------------------------------------

	fmt.Println("For action", action, "fraud type is ", fraud)
	// 封装响应数据
	data := &struct {
		Paper    service.InventoryFinancingPaper
		Msg      string
		Flag     bool
		Username string
		Role     string
		Action   string
		Fraud	 string
	}{
		Paper:    paper,
		Msg:      "",
		Flag:     false,
		Username: username,
		Role:     role,
		Action:   action,
		Fraud:    fraud,
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
	Port := r.FormValue("Port")

	// 获取用户信息
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"].(string)
	role := session.Values["role"].(string)

	msg := ""
	// 调用业务层, 反序列化
	if OrgName != "" {
		msg, _ = app.Fabric.QueryChan(OrgName, Port)
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

// // 根据指定的 Key 查询信息
// func (app *Application) QueryInfo(w http.ResponseWriter, r *http.Request) {
// 	// 获取提交数据
// 	jeweler := r.FormValue("jeweler")
// 	paperNumber := r.FormValue("paperNumber")
// 	action := "QueryPaper"

// 	// 获取用户信息
// 	session, _ := store.Get(r, "mysession")
// 	username := session.Values["username"].(string)
// 	role := session.Values["role"].(string)
// 	// fmt.Println("username: ", username)
// 	// fmt.Println("role: ", role)

// 	// 调用业务层, 反序列化
// 	msg, err := app.Fabric.Querypaper(jeweler, paperNumber)
// 	var paper = service.InventoryFinancingPaper{}
// 	json.Unmarshal(msg, &paper)
// 	fmt.Println(paper)
// 	fmt.Println("paper state: ", paper.State, "\n")

// 	// 封装响应数据
// 	data := &struct {
// 		Paper    service.InventoryFinancingPaper
// 		Msg      string
// 		Flag     bool
// 		Username string
// 		Role     string
// 		Action   string
// 	}{
// 		Paper:    paper,
// 		Msg:      "",
// 		Flag:     false,
// 		Username: username,
// 		Role:     role,
// 		Action:   action,
// 	}
// 	if err != nil {
// 		data.Msg = err.Error()
// 		data.Flag = true
// 	}
// 	// 响应客户端
// 	showView(w, r, "queryResult.html", data)
// }
