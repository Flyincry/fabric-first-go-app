package main

import (
	"fmt"
	"os"

	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/sdkenv"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/service"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/web"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/web/controllers"
)

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

	// sdk setup
	sdk, err := sdkenv.Setup(os.Getenv("GOPATH")+"/src/fabric-first-go-app/config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	fmt.Println(">> 通过链码外部服务设置链码状态......")
	serviceHandler, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	if err != nil {
		fmt.Println()
		os.Exit(-1)
	}

	fmt.Println(">> 设置链码状态完成")

	fmt.Println(">> 启动web服务......")
	app := controllers.Application{
		Fabric: serviceHandler,
	}
	web.WebStart(&app)
	fmt.Println(">> 启动web服务......")
}
