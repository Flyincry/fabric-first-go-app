package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/sdkenv"
)

func Joinchannel(OrgID string, ChannelName string) error {
	orgs := []*sdkenv.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org" + OrgID,
			OrgMspId:      "Org" + OrgID + "MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/fabric-samples/test-network/channel-artifacts/Org" + OrgID + "MSPanchors.tx",
		},
	}

	info := sdkenv.SdkEnvInfo{
		ChannelID:        ChannelName,
		ChannelConfig:    os.Getenv("GOPATH") + "/src/fabric-samples/test-network/channel-artifacts/" + ChannelName + ".tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      "simplecc",
		ChaincodePath:    os.Getenv("GOPATH") + "/src/fabric-first-go-app/newIFC/",
		ChaincodeVersion: "1.0.0",
	}

	_, err := sdkenv.Setup(os.Getenv("GOPATH")+"/src/fabric-first-go-app/config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	fmt.Println(">> 加入通道......")
	for _, org := range info.Orgs {
		// 加入通道
		// Org peers join channel
		if err := org.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com")); err != nil {
			return fmt.Errorf("%s peers failed to JoinChannel: %v", org.OrgName, err)
		}
	}
	fmt.Println(">> 加入通道成功")
	return nil
}

func main() {
	Joinchannel("3", "testchannel")
}
