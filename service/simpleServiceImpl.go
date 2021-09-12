package service

import (
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

func (t *ServiceHandler) Querypaper(jeweler, paperNumber string) (string, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "QueryPaper", Args: [][]byte{[]byte(jeweler), []byte(paperNumber)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}

func (t *ServiceHandler) Apply(paperNumber, jeweler, applyDateTime, financialAmount string) (string, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "Apply", Args: [][]byte{[]byte(paperNumber), []byte(jeweler), []byte(applyDateTime), []byte(financialAmount)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}
