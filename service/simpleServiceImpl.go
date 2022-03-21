package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/shuizhongmose/go-fabric/fabric-first-go-app/db/model"
)

func (t *ServiceHandler) SetInfo(name, num string) (string, error) {

	eventID := "eventSetInfo"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "set", Args: [][]byte{[]byte(name), []byte(num), []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceHandler) Querypaper(jeweler, paperNumber string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "QueryPaper", Args: [][]byte{[]byte(jeweler), []byte(paperNumber)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return response.Payload, nil
}

func (t *ServiceHandler) Apply(paperNumber, jeweler, applyDateTime, financialAmount string) (string, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "Apply", Args: [][]byte{[]byte(paperNumber), []byte(jeweler), []byte(applyDateTime), []byte(financialAmount)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceHandler) Action(paper InventoryFinancingPaper, Action string) ([]byte, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	var comm []string
	switch Action {
	case "Accept", "Payback", "Repurchase":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, now} //4
	case "Supervise", "Default", "QueryPaper", "Reject":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler} //3
	case "Apply":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.JewelerAddr, now, paper.FinancingAmount, paper.PledgeType, paper.PledgeAmount, paper.PledgeApraisedValue} //9
	case "Revise":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.FinancingAmount, now} //5
	case "Evaluate":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.Evaluator, paper.EvalType, paper.EvalQualityProportion, paper.EvalAmount, paper.EvalPrice, now} //9
	case "Receive":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.Bank, now} //5
	case "ReadyRepo":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.Repurchaser, now} //5
	case "PutInStorage":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.Supervisor, paper.StorageAmount, paper.StorageType, paper.StorageAddress, paper.StartDate, paper.EndDate, now} //10
	case "OfferProductInfo":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.Productor, paper.ProductType, paper.ProductAmount, paper.ProductDate, now} //8
	case "OfferLisenceInfo":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.BrandCompany, paper.BrandCompanyAddr, paper.GrantedObject, paper.GrantedStartDate, paper.GrantedEndDate, now} //9
	}
	//fmt.Println("Action comm:" + comm[0] + comm[1] + comm[2] + comm[3] + comm[4] + comm[5] + comm[6] + comm[7] + comm[8] + "\n")
	var response channel.Response
	var err error
	switch len(comm) {
	case 3:
		response, err = t.Client.Execute(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2])}})
	case 4:
		response, err = t.Client.Execute(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3])}})
	case 5:
		response, err = t.Client.Execute(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3]), []byte(comm[4])}})
	case 8:
		response, err = t.Client.Execute(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3]), []byte(comm[4]), []byte(comm[5]), []byte(comm[6]), []byte(comm[7])}})
	case 9:
		response, err = t.Client.Execute(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3]), []byte(comm[4]), []byte(comm[5]), []byte(comm[6]), []byte(comm[7]), []byte(comm[8])}})
	case 10:
		response, err = t.Client.Execute(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3]), []byte(comm[4]), []byte(comm[5]), []byte(comm[6]), []byte(comm[7]), []byte(comm[8]), []byte(comm[9])}})
	}

	if err != nil {
		return []byte{0x00}, err
	}

	return response.Payload, nil
}

func (t *ServiceHandler) CreateChan(ChannelName string) {
	cmd := exec.Command("./network.sh", "createChannel", "-c", ChannelName)
	cmd.Dir = "/root/workspace/src/fabric-samples/test-network"
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Print("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

func (t *ServiceHandler) CreateOrg(OrgID string) error {
	cmd := exec.Command("./runme.sh", OrgID)
	cmd.Dir = "/root/workspace/src/fabric-samples/test-network/addOrg3"
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Print("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
	return err
}

func (t *ServiceHandler) QueryChan(Name string) (Msg string, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获异常:", err)
		}
	}()
	// Get user orgid and port from database
	u := model.User{
		Name: Name,
	}
	u1, err := u.GetUserByName()
	if err != nil {
		return "Unexist OrgName.", err
	}
	OrgName := u1.OrgID
	orgid, _ := strconv.Atoi(u1.OrgID)
	Port := strconv.Itoa(orgid*2000 + 5051)
	// Setting environment variables
	os.Setenv("FABRIC_CFG_PATH", "/root/workspace/src/fabric-samples/config")
	os.Setenv("CORE_PEER_TLS_ENABLED", "true")
	os.Setenv("CORE_PEER_LOCALMSPID", "Org"+OrgName+"MSP")
	os.Setenv("CORE_PEER_TLS_ROOTCERT_FILE", "/root/workspace/src/fabric-samples/test-network/organizations/peerOrganizations/org"+OrgName+".example.com/peers/peer0.org"+OrgName+".example.com/tls/ca.crt")
	os.Setenv("CORE_PEER_MSPCONFIGPATH", "/root/workspace/src/fabric-samples/test-network/organizations/peerOrganizations/org"+OrgName+".example.com/users/Admin@org"+OrgName+".example.com/msp")
	os.Setenv("CORE_PEER_ADDRESS", "localhost:"+Port)

	cmd := exec.Command("peer", "channel", "list")
	cmd.Dir = "/root/workspace/src/fabric-samples/test-network"

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		return "Blockchain Network Connection Error.", err
	}
	fmt.Printf("combined out:\n%s\n", string(out))
	location := strings.IndexAny(string(out), "joined:") + 130
	rawStrSlice := []byte(string(out))
	res := string(rawStrSlice[location:])
	return res, err
}

func (t *ServiceHandler) JoinChan(Name, ChannelName string) (Msg string, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获异常:", err)
		}
	}()
	// Get user orgid and port from database
	u := model.User{
		Name: Name,
	}
	u1, err := u.GetUserByName()
	if err != nil {
		return "Unexist OrgName.", err
	}
	OrgName := u1.OrgID
	orgid, _ := strconv.Atoi(u1.OrgID)
	Port := strconv.Itoa(orgid*2000 + 5051)
	// Setting environment variables
	os.Setenv("FABRIC_CFG_PATH", "/root/workspace/src/fabric-samples/config")
	os.Setenv("CORE_PEER_TLS_ENABLED", "true")
	os.Setenv("CORE_PEER_LOCALMSPID", "Org"+OrgName+"MSP")
	os.Setenv("CORE_PEER_TLS_ROOTCERT_FILE", "/root/workspace/src/fabric-samples/test-network/organizations/peerOrganizations/org"+OrgName+".example.com/peers/peer0.org"+OrgName+".example.com/tls/ca.crt")
	os.Setenv("CORE_PEER_MSPCONFIGPATH", "/root/workspace/src/fabric-samples/test-network/organizations/peerOrganizations/org"+OrgName+".example.com/users/Admin@org"+OrgName+".example.com/msp")
	os.Setenv("CORE_PEER_ADDRESS", "localhost:"+Port)

	cmd := exec.Command("peer", "channel", "join", "-b", "./channel-artifacts/"+ChannelName+".block")
	cmd.Dir = "/root/workspace/src/fabric-samples/test-network"

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		return "Fail to join channel. Possible reason: unexisted channel.", err
	}
	fmt.Printf("combined out:\n%s\n", string(out))
	location := strings.IndexAny(string(out), "executeJoin") + 162
	rawStrSlice := []byte(string(out))
	res := string(rawStrSlice[location:])
	return res, err
}
