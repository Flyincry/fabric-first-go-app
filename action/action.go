package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/sdkenv"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/service"
)

//珠宝商
var Action = flag.String("action", "", "动作名")
var Jeweler = flag.String("jeweler", "", "珠宝商")
var PaperNumber = flag.String("paperNumber", "", "申请文书ID")
var FinancingAmount = flag.String("financingAmount", "", "融资金额")

//生产者
var Productor = flag.String("productor", "", "生产商")
var ProductType = flag.String("productType", "", "货品种类")
var ProductAmount = flag.String("productAmount", "", "货品数量")
var ProductDate = flag.String("productDate", "", "货品生产日期")

//品牌方
var BrandCompany = flag.String("brandCompany", "", "品牌方")
var GrantedObject = flag.String("grantedObject", "", "授权对象")
var GrantedStartDate = flag.String("grantedStartDate", "", "授权开始日期")
var GrantedEndDate = flag.String("grantedEndDate", "", "授权结束日期")

//银行
var Bank = flag.String("bank", "", "银行")

//评估鉴定方
var Evaluator = flag.String("evaluator", "", "评估者")
var EvalType = flag.String("evalType", "", "评估种类")
var EvalQualityProportion = flag.String("evalQualityProportion", "", "评估质量（质检合格比例）")
var EvalAmount = flag.String("evalAmount", "", "评估价值")

//仓库监管方
var Supervisor = flag.String("supervisor", "", "监管者")
var StorageAmount = flag.String("storageAmount", "", "仓库货品总量")
var StorageType = flag.String("storageType", "", "货品种类")
var StorageAddress = flag.String("storageAddress", "", "存储地址")
var EndDate = flag.String("endDate", "", "期限")

//回购方
var Repurchaser = flag.String("repurchaser", "", "回购方")

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

	flag.Parse()
	now := time.Now().Format("2006-01-02 15:04:05")
	var comm []string
	switch *Action {
	case "Accept", "Payback", "Repurchase":
		comm = []string{*Action, *PaperNumber, *Jeweler, now} //4
	case "Supervise", "Default", "QueryPaper", "Reject":
		comm = []string{*Action, *PaperNumber, *Jeweler} //3
	case "Apply", "Revise":
		comm = []string{*Action, *PaperNumber, *Jeweler, *FinancingAmount, now} //5
	case "Evaluate":
		comm = []string{*Action, *PaperNumber, *Jeweler, *Evaluator, *EvalType, *EvalQualityProportion, *EvalAmount, now} //8
	case "Receive":
		comm = []string{*Action, *PaperNumber, *Jeweler, *Bank, now} //5
	case "ReadyRepo":
		comm = []string{*Action, *PaperNumber, *Jeweler, *Repurchaser, now} //5
	case "PutInStorage":
		comm = []string{*Action, *PaperNumber, *Jeweler, *Supervisor, *StorageAmount, *StorageType, *StorageAddress, *EndDate, now} //9
	case "OfferProductInfo":
		comm = []string{*Action, *PaperNumber, *Jeweler, *Productor, *ProductType, *ProductAmount, *ProductDate, now} //8
	case "OfferLisenceInfo":
		comm = []string{*Action, *PaperNumber, *Jeweler, *BrandCompany, *GrantedObject, *GrantedStartDate, *GrantedEndDate, now} //9
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
	}

	if err != nil {
		fmt.Println(">>failed : %v", err)
	}

	fmt.Println("<--- 添加信息　--->：", string(response.Payload))
}
