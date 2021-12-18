package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
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
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

func (t *ServiceHandler) Apply(paperNumber, jeweler, applyDateTime, financialAmount string) (string, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "Apply", Args: [][]byte{[]byte(paperNumber), []byte(jeweler), []byte(applyDateTime), []byte(financialAmount)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceHandler) Action(paper InventoryFinancingPaper, Action string) (string, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	var comm []string
	switch Action {
	case "Accept", "Payback", "Repurchase":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, now} //4
	case "Supervise", "Default", "QueryPaper", "Reject":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler} //3
	case "Apply", "Revise":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.FinancingAmount, now} //5
	case "Evaluate":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.Evaluator, paper.EvalType, paper.EvalQualityProportion, paper.EvalAmount, now} //8
	case "Receive":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.Bank, now} //5
	case "ReadyRepo":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.Repurchaser, now} //5
	case "PutInStorage":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.Supervisor, paper.StorageAmount, paper.StorageType, paper.StorageAddress, paper.EndDate, now} //9
	case "OfferProductInfo":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.Productor, paper.ProductType, paper.ProductAmount, paper.ProductDate, now} //8
	case "OfferLisenceInfo":
		comm = []string{Action, paper.PaperNumber, paper.Jeweler, paper.BrandCompany, paper.GrantedObject, paper.GrantedStartDate, paper.GrantedEndDate, now} //9
	}

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
	}

	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}

func (t *ServiceHandler) CreateChan(ChannelName string) {
	cmd := exec.Command("./network.sh", "createChannel", "-c", ChannelName)
	cmd.Dir = "/root/workspace/src/fabric-samples/test-network"
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

func (t *ServiceHandler) CreateOrg(OrgID string) error {
	cmd := exec.Command("./runme.sh", OrgID)
	cmd.Dir = "/root/workspace/src/fabric-samples/test-network/addOrg3"
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
	return err
}

func (t *ServiceHandler) QueryChan(OrgName, Port string) (Msg string) {
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
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
	return string(out)
}

// func (t *ServiceHandler) Action2(action, Jeweler, PaperNumber, FinancingAmount, ApplyDateTime, ReviseDateTime, AcceptDateTime, ReadyDateTime, EvalDateTime, ReceiveDateTime, EndDate, PaidbackDateTime, RepurchaseDateTime, Bank, Evaluator, Repurchaser, Supervisor string) (string, error) {
// 	var comm []string
// 	switch action {
// 	case "Accept":
// 		comm = []string{action, Jeweler, PaperNumber, AcceptDateTime}
// 	case "Apply":
// 		comm = []string{action, PaperNumber, Jeweler, ApplyDateTime, FinancingAmount}
// 	case "Default", "QueryPaper", "Reject":
// 		comm = []string{action, Jeweler, PaperNumber}
// 	case "Evaluate":
// 		comm = []string{action, Jeweler, PaperNumber, Evaluator, EvalDateTime}
// 	case "Payback":
// 		comm = []string{action, Jeweler, PaperNumber, PaidbackDateTime}
// 	case "ReadyRepo":
// 		comm = []string{action, Jeweler, PaperNumber, Repurchaser, ReadyDateTime}
// 	case "Receive":
// 		comm = []string{action, Jeweler, Bank, PaperNumber, ReceiveDateTime}
// 	case "Repurchase":
// 		comm = []string{action, Jeweler, PaperNumber, RepurchaseDateTime}
// 	case "Revise":
// 		comm = []string{action, Jeweler, PaperNumber, ReviseDateTime, FinancingAmount}
// 	case "Supervise":
// 		comm = []string{action, Jeweler, Supervisor, EndDate, PaperNumber}
// 	}

// 	var response channel.Response
// 	var err error
// 	switch len(comm) {
// 	case 3:
// 		response, err = t.Client.Execute(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2])}})
// 	case 4:
// 		response, err = t.Client.Execute(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3])}})
// 	case 5:
// 		response, err = t.Client.Execute(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: comm[0], Args: [][]byte{[]byte(comm[1]), []byte(comm[2]), []byte(comm[3]), []byte(comm[4])}})
// 	}

// 	if err != nil {
// 		return "", err
// 	}

// 	return string(response.TransactionID), nil
// }
