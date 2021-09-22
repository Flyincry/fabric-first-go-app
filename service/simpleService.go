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

type InventoryFinancingPaper struct {
	PaperNumber        string `json:"paperNumber"`
	Jeweler            string `json:"jeweler"`
	ApplyDateTime      string `json:"applyDateTime"`
	ReviseDateTime     string `json:"reviseDateTime"`
	AcceptDateTime     string `json:"acceptDateTime"`
	ReadyDateTime      string `json:"readyDateTime"`
	EvalDateTime       string `json:"evalDateTime"`
	ReceiveDateTime    string `json:"receiveDateTime"`
	EndDate            string `json:"endDateTime"`
	PaidbackDateTime   string `json:"paidBackDateTime"`
	RepurchaseDateTime string `json:"RepurchaseDateTime"`
	FinancingAmount    int    `json:"financingAmount"`
	Dealer             string `json:"dealer"`
	State              int    `json:"currentState"`
	Bank               string `json:"bank"`
	Evaluator          string `json:"evaluator"`
	Repurchaser        string `json:"repurchaser"`
	Supervisor         string `json:"supervisor"`
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
