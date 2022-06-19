package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"bytes"
	"net"
	"log"
	"encoding/json"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"errors"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/sdkenv"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/service"
)
 
type State uint

// InventoryFinancingPaper 定义了一个珠宝存货融资流程
type InventoryFinancingPaper struct {
	//珠宝商发起融资申请
	PaperNumber         string `json:"paperNumber"`         //融资申请编号
	Jeweler             string `json:"jeweler"`             //融资申请珠宝商
	JewelerAddr         string `json:"jewelerAddr"`         //融资申请珠宝商门店地址
	ApplyDateTime       string `json:"applyDateTime"`       //提交申请时间（web端自动生成）
	FinancingAmount     string `json:"financingAmount"`     //融资金额 **** int->string
	PledgeType          string `json:"pledgeType"`          //质押的货品类别****** new
	PledgeAmount        string `json:"pledgeAmount"`        //质押货品数量****** new
	PledgeApraisedValue string `json:"pledgeApraisedValue"` //质押货品预估价值****** new
	//生产者提供了生产信息上链
	Productor             string `json:"productor"`             //生产商
	ProductType           string `json:"productType"`           //货品种类
	ProductAmount         string `json:"productAmount"`         //货品数量 int->string
	ProductDate           string `json:"productDate"`           //货品生产日期
	ProductInfoUpdateTime string `json:"productInfoUpdateTime"` //货品信息更新日期（web端自动生成）
	//品牌方提供授权信息上链
	BrandCompany          string `json:"brandCompany"`         //品牌方
	BrandCompanyAddr      string `json:"brandCompanyAddr"`     //品牌方地址
	GrantedObject         string `json:"grantedObject"`        //授权对象
	GrantedStartDate      string `json:"grantedStartDate"`      //授权开始日期
	GrantedEndDate        string `json:"grantedEndDate"`        //授权结束日期
	GrantedInfoUpdateTime string `json:"grantedInfoUpdateTime"` //授权信息更新日期（web端自动生成）
	//银行收到融资申请
	Bank            string `json:"bank"`
	ReceiveDateTime string `json:"receiveDateTime"` //收到融资申请时间（web端自动生成）
	//评估鉴定方提供鉴定信息
	Evaluator             string `json:"evaluator"`
	EvalDateTime          string `json:"evalDateTime"`          //鉴定时间（web端自动生成）
	EvalType              string `json:"evalType"`              //评估种类
	EvalQualityProportion string `json:"evalQualityProportion"` //评估质量（质检合格比例）
	EvalAmount            string `json:"evalAmount"`            //评估产品数量 int->string
	EvalPrice             string `json:"evalPrice"`             //评估价格
	//仓库监管方提供仓单信息
	Supervisor        string `json:"supervisor"`
	StorageAmount     string `json:"storageAmount"`     //仓库货品总量
	StorageType       string `json:"storageType"`       //货品种类
	StorageAddress    string `json:"storageAddress"`    //存储地址
	StartDate         string `json:"startDate"`         //融资开始时间******* new
	EndDate           string `json:"endDate"`           //融资终止时间
	StorageInfoUpdate string `json:"storageInfoUpdate"` //出具仓单的时间（web端自动生成）
	//回购方准备好可以后续回购
	Repurchaser   string `json:"repurchaser"`
	ReadyDateTime string `json:"readyDateTime"`
	//银行接受
	AcceptDateTime string `json:"acceptedDateTime"` //银行接受时间（web端自动生成）
	//珠宝商回购
	PaidbackDateTime string `json:"paidBackDateTime"` //珠宝商回购时间（web端自动生成）
	//回购方回购
	RepurchaseDateTime string `json:"repurchaseDateTime"` //回购方回购时间（web端自动生成）
	state              State  `metadata:"currentState"`
	prevstate          State  `metadata:"prevState,optional"`
	class              string `metadata:"class"`
	key                string `metadata:"key"`
	//珠宝商重写
	ReviseDateTime string `json:"reviseDateTime"`
}


var Action string
var action_from_flag = flag.String("Action", "", "The action from airflow")

//珠宝商
var Jeweler = string("jeweler")
var PaperNumber = string("paperNumber")
var FinancingAmount = string("financingAmount")
var JewelerAddr = string("jewelerAddr")
var PledgeType = string("pledgeType")
var PledgeAmount = string("pledgeAmount")
var PledgeApraisedValue = string("pledgeApraisedValue")

//生产者
var Productor = string("productor")
var ProductType = string("productType")
var ProductAmount = string("productAmount")
var ProductDate = string("productDate")

//品牌方
var BrandCompany = string("brandCompany")
var BrandCompanyAddr = string("brandCompanyAddr")
var GrantedObject = string("grantedObject")
var GrantedStartDate = string("grantedStartDate")
var GrantedEndDate = string("grantedEndDate")

//银行
var Bank = string("bank")

//评估鉴定方
var Evaluator = string("evaluator")
var EvalType = string("evalType")
var EvalQualityProportion = string("evalQualityProportion")
var EvalAmount = string("evalAmount")
var EvalPrice = string("evalPrice")

//仓库监管方
var Supervisor = string("supervisor")
var StorageAmount = string("storageAmount")
var StorageType = string("storageType")
var StorageAddress = string("storageAddress")
var StartDate = string("startDate")
var EndDate = string("endDate")

//回购方
var Repurchaser = string("repurchaser")

// 解析私匙
func parsingRsaPrivateKey(file string) (*rsa.PrivateKey, error) {
	// 读取私匙
	priByte,err := ioutil.ReadFile(file)
	if err != nil {
	  return nil,err
	}
	// pem解码
	b,_ := pem.Decode(priByte)
	if b == nil {
	  return nil,errors.New("error private key")
	}
	// der加密，返回一个私匙对象
	prikey,err := x509.ParsePKCS1PrivateKey(b.Bytes)
	if err != nil {
	  return nil,err
	}
	return prikey,nil
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

func main() {

	orgs := []*sdkenv.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/fabric-samples/test-network/channel-artifacts/Org1MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/fabric-samples/test-network/channel-artifacts/Org2MSPanchors.tx",
		},
	}

	// init sdk env info
	info := sdkenv.SdkEnvInfo{
		ChannelID:        "testchannel",
		ChannelConfig:    os.Getenv("GOPATH") + "/src/fabric-samples/test-network/channel-artifacts/testchannel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      "simplecc",
		ChaincodePath:    os.Getenv("GOPATH") + "/src/fabric-first-go-app/IFC/",
		ChaincodeVersion: "1.0.0",
	}

	sdk, err := sdkenv.Setup(os.Getenv("GOPATH")+"/src/fabric-first-go-app/config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}
	fmt.Println(">> setup successful......")

	fmt.Println(">> 通过链码外部服务设置链码状态......")
	App, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	if err != nil {
		fmt.Println()
		os.Exit(-1)
	}
	fmt.Println(">> 设置链码状态完成")

	// 解析airflow传的flag参数
	flag.Parse()
	// parse flag 之后已经可以调用action_from_flag了

	//加入client-server之前是通过flag解析参数，参数写死在airflow dags里面
	// flag.Parse()

	// ------------------------------
	// add a server
	// 要根据action_from_flag来看开哪个端口，不然数据可能会传错
	// offerproductinfo 8900
	// offerliscenceinfo 8901
	// evaluate 8902
	// putinstorage 8903
	// readyrepo 8904
	// others 8888
	
	var PORT int
	// PORT = 8888
	if *action_from_flag == "OfferProductInfo"{
		PORT = 8900
	}else if *action_from_flag == "OfferLisenceInfo"{
		PORT = 8901
	}else if *action_from_flag == "Evaluate"{
		PORT = 8902
	}else if *action_from_flag == "PutInStorage"{
		PORT = 8903
	}else if *action_from_flag == "ReadyRepo"{
		PORT = 8904
	}else{
		PORT = 8888
	}
	fmt.Println("The action from flag is :", *action_from_flag, "port number is :", PORT)

	address := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"), // 把字符串IP地址转换为net.IP类型
		Port: PORT,
	}
    listener, err := net.ListenTCP("tcp4", &address) // 创建TCP4服务器端监听器
    if err != nil {
        log.Fatal(err) // Println + os.Exit(1)
    }

	fmt.Println("开始监听")
    conn, err := listener.AcceptTCP()
    if err != nil {
        log.Fatal(err) // 错误直接退出
    }

    json_message := make([]byte, 1024)
    n, err := conn.Read(json_message)
	fmt.Println(n)

    fmt.Println("before print")
	fmt.Println(string(json_message[:n])+"\n")
	fmt.Println("after print")

	// rsa 解密
	priKey,_ := parsingRsaPrivateKey("pri.key") // 解密私匙
	decryData,_ := rsaPrivateKeyDecrypt(json_message[:n],priKey) // 解密数据
	fmt.Printf("%s\n",decryData)

	// 以结构体的形式输出paper
    var out bytes.Buffer
	json.Indent(&out, decryData, "", "\t")
	fmt.Printf("结构体paper=%v\n", out.String())

	// 解决TCP粘包问题
	time.Sleep(time.Duration(2)*time.Second)

    action_message := make([]byte, 256)
    nn, err1 := conn.Read(action_message)
    fmt.Println(string(action_message[:nn]))
    Action = string(action_message[:nn])
    
    fmt.Println("action in the loop: " + Action)
    if err1 != nil {
		fmt.Println("err1 happens.")
		fmt.Println(err1)
	}

    var paper InventoryFinancingPaper 
	// json转为结构体
    err2 := json.Unmarshal(json_message[:n], &paper)
    fmt.Println(paper.PaperNumber)
	fmt.Println("输出品牌公司试试")
	fmt.Println(paper.BrandCompany)
    
	if err2 != nil {
		fmt.Println("err2 happens.")
		fmt.Println(err2)
	}
    
    fmt.Println("remote address:", conn.RemoteAddr())
    fmt.Println("收到参数了")
	fmt.Println("1Action is :", Action)
	// 关闭server
	conn.Close()

	if *action_from_flag != Action{
		fmt.Println("The action from client is different from the action in Airflow!")
	}

	//珠宝商
	Jeweler = paper.Jeweler
	PaperNumber = paper.PaperNumber
	FinancingAmount = paper.FinancingAmount
	JewelerAddr = paper.JewelerAddr
	PledgeType = paper.PledgeType
	PledgeAmount = paper.PledgeAmount
	PledgeApraisedValue = paper.PledgeApraisedValue

	//生产者
	Productor = paper.Productor
	ProductType = paper.ProductType
	ProductAmount = paper.ProductAmount
	ProductDate = paper.ProductDate

	//品牌方
	BrandCompany = paper.BrandCompany
	BrandCompanyAddr = paper.BrandCompanyAddr
	GrantedObject = paper.GrantedObject
	GrantedStartDate = paper.GrantedStartDate
	GrantedEndDate = paper.GrantedEndDate

	fmt.Println("print out the information of brand company:", BrandCompany, BrandCompanyAddr, GrantedObject, GrantedStartDate, GrantedEndDate)

	//银行
	Bank = paper.Bank

	//评估鉴定方
	Evaluator = paper.Evaluator
	EvalType = paper.EvalType
	EvalQualityProportion = paper.EvalQualityProportion
	EvalAmount = paper.EvalAmount
	EvalPrice = paper.EvalPrice

	//仓库监管方
	Supervisor = paper.Supervisor
	StorageAmount = paper.StorageAddress
	StorageType = paper.StorageType
	StorageAddress = paper.StorageAddress
	StartDate = paper.StartDate
	EndDate = paper.EndDate

	//回购方
	Repurchaser = paper.Repurchaser
        
	fmt.Println("2Action is :", Action)

	//------------------------------



	now := time.Now().Format("2006-01-02 15:04:05")
	var comm []string
	switch Action {
	case "Accept", "Payback", "Repurchase":
		comm = []string{Action, PaperNumber, Jeweler, now} //4
	case "Supervise", "Default", "QueryPaper", "Reject":
		comm = []string{Action, PaperNumber, Jeweler} //3
	case "Apply":
		comm = []string{Action, PaperNumber, Jeweler, FinancingAmount, JewelerAddr, now, FinancingAmount, PledgeType, PledgeAmount, PledgeApraisedValue} //9
	case "Revise":
		comm = []string{Action, PaperNumber, Jeweler, FinancingAmount, now} //5
	case "Evaluate":
		comm = []string{Action, PaperNumber, Jeweler, Evaluator, EvalType, EvalQualityProportion, EvalAmount, EvalPrice, now} //9
	case "Receive":
		comm = []string{Action, PaperNumber, Jeweler, Bank, now} //5
	case "ReadyRepo":
		comm = []string{Action, PaperNumber, Jeweler, Repurchaser, now} //5
	case "PutInStorage":
		comm = []string{Action, PaperNumber, Jeweler, Supervisor, StorageAmount, StorageType, StorageAddress, StartDate, EndDate, now} //10
	case "OfferProductInfo":
		comm = []string{Action, PaperNumber, Jeweler, Productor, ProductType, ProductAmount, ProductDate, now} //8
	case "OfferLisenceInfo":
		comm = []string{Action, PaperNumber, Jeweler, BrandCompany, BrandCompanyAddr, GrantedObject, GrantedStartDate, GrantedEndDate, now} //9
	}

	var response channel.Response
	switch len(comm) {
	case 3:
		response, err = App.Client.Execute(channel.Request{ChaincodeID: App.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2])}})
	case 4:
		response, err = App.Client.Execute(channel.Request{ChaincodeID: App.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3])}})
	case 5:
		response, err = App.Client.Execute(channel.Request{ChaincodeID: App.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3]), []byte(comm[4])}})
	case 8:
		response, err = App.Client.Execute(channel.Request{ChaincodeID: App.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3]), []byte(comm[4]), []byte(comm[5]), []byte(comm[6]), []byte(comm[7])}})
	case 9:
		response, err = App.Client.Execute(channel.Request{ChaincodeID: App.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3]), []byte(comm[4]), []byte(comm[5]), []byte(comm[6]), []byte(comm[7]), []byte(comm[8])}})
	case 10:
		response, err = App.Client.Execute(channel.Request{ChaincodeID: App.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3]), []byte(comm[4]), []byte(comm[5]), []byte(comm[6]), []byte(comm[7]), []byte(comm[8]), []byte(comm[9])}})
	}

	if err != nil {
		fmt.Println(">>failed : %v", err)
	}

	fmt.Println("<--- 添加信息　--->：", string(response.Payload))
}

// 添加flag
// 改dags
// 走完流程
