package service

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/sdkenv"
)

type ServiceHandler struct {
	ChaincodeID string
	Client      *channel.Client
}

type State uint

type InventoryFinancingPaper struct {
	//珠宝商发起融资申请
	PaperNumber         string `json:"paperNumber"`         //融资申请编号
	Jeweler             string `json:"jeweler"`             //融资申请珠宝商
	JewelerAddr         string `json:"jewelerAddr"`         //融资申请珠宝商门店地址****** new
	ApplyDateTime       string `json:"applyDateTime"`       //提交申请时间（web端自动生成）
	FinancingAmount     string `json:"financingAmount"`     //融资金额
	PledgeType          string `json:"pledgeType"`          //质押的货品类别****** new
	PledgeAmount        string `json:"pledgeAmount"`        //质押货品数量****** new
	PledgeApraisedValue string `json:"pledgeApraisedValue"` //质押货品预估价值****** new
	//生产者提供了生产信息上链
	Productor             string `json:"productor"`             //生产商
	ProductType           string `json:"productType"`           //货品种类
	ProductAmount         string `json:"productAmount"`         //货品数量
	ProductDate           string `json:"productDate"`           //货品生产日期
	ProductInfoUpdateTime string `json:"productInfoUpdateTime"` //货品信息更新日期（web端自动生成）
	//品牌方提供授权信息上链
	BrandCompany          string `json:"brandCompany"`          //品牌方
	BrandCompanyAddr      string `json:"brandCompanyAddr"`      //品牌方地址****** new
	GrantedObject         string `json:"grantedObject"`         //授权对象
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
	EvalAmount            string `json:"evalAmount"`            //评估数量
	EvalPrice             string `json:"evalPrice"`             //评估价格****** new
	//仓库监管方提供仓单信息
	Supervisor        string `json:"supervisor"`
	StorageAmount     string `json:"storageAmount"`     //仓库货品总量
	StorageType       string `json:"storageType"`       //货品种类
	StorageAddress    string `json:"storageAddress"`    //存储地址
	StartDate         string `json:"startDate"`         //融资开始时间******* new
	EndDate           string `json:"endDate"`           //融资结束时间
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
	State              State  `json:"currentState"`
	class              string `metadata:"class"`
	key                string `metadata:"key"`
	//珠宝商重写
	ReviseDateTime string `json:"reviseDateTime"`
}

func InitService(chaincodeID, channelID string, org *sdkenv.OrgInfo, sdk *fabsdk.FabricSDK) (*ServiceHandler, error) {
	handler := &ServiceHandler{
		ChaincodeID: chaincodeID,
	}
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new channel client: %s", err)
	}
	handler.Client = client
	return handler, nil
}

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}
