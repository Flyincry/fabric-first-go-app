/*
 * SPDX-License-Identifier: Apache-2.0
 */

package inventoryfinancingpaper

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract chaincode that defines
// the business logic for managing inventory paper

type Contract struct {
	contractapi.Contract
}

// Init does nothing
func (c *Contract) Init() {
	fmt.Println("Init")
}

// InitLedger adds a base set of InventoryFinancingPapers to the ledger
func (c *Contract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	papers := []InventoryFinancingPaper{
		InventoryFinancingPaper{
			PaperNumber:           "111",
			Jeweler:               "宝琳国金",
			JewelerAddr:           "深圳市罗湖区水贝水田二街3号宝琳国金大厦2栋",
			ApplyDateTime:         "202109",
			FinancingAmount:       "250,000,000",
			PledgeType:            "K金",
			PledgeAmount:          "1000kg",
			PledgeApraisedValue:   "261,000,000",
			Productor:             "深圳罗玛星珠宝首饰有限公司",
			ProductType:           "K金",
			ProductAmount:         "1000kg",
			ProductDate:           "202106",
			ProductInfoUpdateTime: "202111",
			BrandCompany:          "周大生",
			BrandCompanyAddr:      "深圳市罗湖区翠竹街道布心路3033号水贝壹号A座19-23层",
			GrantedObject:         "宝琳国金",
			GrantedStartDate:      "202001",
			GrantedEndDate:        "202203",
			GrantedInfoUpdateTime: "202111",
			Bank:                  "深圳发展银行",
			ReceiveDateTime:       "202111",
			Evaluator:             "深圳国艺珠宝艺术品资产评估有限公司",
			EvalDateTime:          "202112",
			EvalType:              "K金",
			EvalQualityProportion: "99%",
			EvalAmount:            "1000kg",
			EvalPrice:             "260,000,000",
			Supervisor:            "宝时云仓",
			StorageAmount:         "1000kg",
			StorageType:           "K金",
			StorageAddress:        "松江一仓-上海松江区泗砖路351号7号楼",
			StartDate:             "202103",
			EndDate:               "202203",
			StorageInfoUpdate:     "202111",
			Repurchaser:           "千禧之星",
			ReadyDateTime:         "",
			AcceptDateTime:        "202111",
			PaidbackDateTime:      "",
			RepurchaseDateTime:    "202203",
			state:                 PAIDBACK,
			prevstate:             SUPERVISING,
			class:                 "",
			key:                   "",
			ReviseDateTime:        "",
		},
	}

	for i, paper := range papers {
		paperAsBytes, _ := json.Marshal(paper)
		err := ctx.GetStub().PutState("InventoryFinancingPaper"+strconv.Itoa(i), paperAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// Apply creates a new inventory paper and stores it in the world state.
func (c *Contract) Apply(ctx TransactionContextInterface, paperNumber string, jeweler string, jewelerAddr string, applyDateTime string, financingAmount string, pledgeType string, pledgeAmount string, pledgeApraisedValue string) (*InventoryFinancingPaper, error) {
	paper := InventoryFinancingPaper{PaperNumber: paperNumber, Jeweler: jeweler, JewelerAddr: jewelerAddr, FinancingAmount: financingAmount, PledgeType: pledgeType, PledgeAmount: pledgeAmount, PledgeApraisedValue: pledgeApraisedValue, ApplyDateTime: applyDateTime}

	paper.SetApplied()
	paper.LogPrevState()

	err := ctx.GetPaperList().AddPaper(&paper)

	if err != nil {
		return nil, err
	}

	fmt.Printf("The jeweler %q  has applied for a new inventory financingp paper %q,the financing amount is %v.\n Current State is %q", jeweler, paperNumber, financingAmount, paper.GetState())
	return &paper, nil
}

// OfferProductInfo means the prodcutor offer the production infomation and  stores it in the world state.
func (c *Contract) OfferProductInfo(ctx TransactionContextInterface, paperNumber string, jeweler string, productor string, productType string, productAmount string, productDate string, productInfoUpdateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)
	if err != nil {
		return nil, err
	}
	if paper.IsApplied() {
		if paper.GetProductor() == "" {
			paper.SetProductor(productor)
			paper.ProductType = productType
			paper.ProductAmount = productAmount
			paper.ProductDate = productDate
			paper.ProductInfoUpdateTime = productInfoUpdateTime

		}
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The productor %q has offed productInfo (productType: %q,productAmount %q, productDate %q, productInfoUpdateTime %q) of the inventory financing paper %q:%q.\n Current State is %q", paper.GetProductor(), productType, productAmount, productDate, productInfoUpdateTime, jeweler, paperNumber, paper.GetState())
	return paper, nil
}

// OfferLisenceInfo means the brand company offer the lisence infomation and  stores it in the world state.
func (c *Contract) OfferLisenceInfo(ctx TransactionContextInterface, paperNumber string, jeweler string, brandCompany string, brandCompanyAddr string, grantedObject string, grantedStartDate string, grantedEndDate string, grantedInfoUpdateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)
	if err != nil {
		return nil, err
	}
	if paper.IsApplied() {
		if paper.GetBrandCompany() == "" {
			paper.SetBrandCompany(brandCompany)
			paper.BrandCompanyAddr = brandCompanyAddr
			paper.GrantedObject = grantedObject
			paper.GrantedStartDate = grantedStartDate
			paper.GrantedEndDate = grantedEndDate
			paper.GrantedInfoUpdateTime = grantedInfoUpdateTime
		}
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The brand company %q has offered LisenceInfo(grantedObject %q, grantedStartDate %q, grantedEndDate %q, grantedInfoUpdateTime %q) of the inventory financing paper %q:%q.\n Current State is %q", paper.GetBrandCompany(), grantedObject, grantedStartDate, grantedEndDate, grantedInfoUpdateTime, jeweler, paperNumber, paper.GetState())
	return paper, nil
}

// Receive updates a inventory paper to be in received status and sets the next dealer
func (c *Contract) Receive(ctx TransactionContextInterface, paperNumber string, jeweler string, bank string, receiveDateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)
	if err != nil {
		return nil, err
	}

	if paper.GetBank() == "" {
		paper.SetBank(bank)
	}

	if paper.GetProductor() != "" && paper.GetBrandCompany() != "" {
		paper.SetReceived()
		paper.ReceiveDateTime = receiveDateTime
	}

	if !paper.IsReceived() {
		return nil, fmt.Errorf("inventory paper %s:%s is not received by bank. Current state = %s", jeweler, paperNumber, paper.GetState())
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The bank %q has received the inventory financing paper %q from jeweler %q,\nCurrent State is %q", paper.GetBank(), paperNumber, jeweler, paper.GetState())
	return paper, nil
}

//Evaluate updates a inventory paper to be evaluated
func (c *Contract) Evaluate(ctx TransactionContextInterface, paperNumber string, jeweler string, evaluator string, evalType string, evalQualityProportion string, evalAmount string, evalPrice string, evalDateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)
	if err != nil {
		return nil, err
	}
	if paper.IsReceived() {
		if paper.GetEvaluator() == "" {
			paper.SetEvaluator(evaluator)
			paper.EvalType = evalType
			paper.EvalQualityProportion = evalQualityProportion
			paper.EvalAmount = evalAmount
			paper.EvalDateTime = evalDateTime
			paper.EvalPrice = evalPrice
		}

	}

	if !paper.IsReceived() {
		return nil, fmt.Errorf("inventory paper %s:%s is not received by bank. Current state = %s", jeweler, paperNumber, paper.GetState())
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The evluator %q has evaluated the inventory financing paper %q:%q.\n The evalType is %q, evalQualityProportion %q, evalAmount %q, evalDateTime %s.The Current State is %q", paper.GetEvaluator(), jeweler, paperNumber, evalType, evalQualityProportion, evalAmount, evalDateTime, paper.GetState())
	return paper, nil
}

//ReadyRepo updates the repurchaser to be ready for Repo
func (c *Contract) ReadyRepo(ctx TransactionContextInterface, paperNumber string, jeweler string, repurchaser string, readyDateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)
	if err != nil {
		return nil, err
	}

	if paper.IsReceived() {
		if paper.GetRepurchaser() == "" {
			paper.SetRepurchaser(repurchaser)
			paper.ReadyDateTime = readyDateTime
		}

	}

	if !paper.IsReceived() {
		return nil, fmt.Errorf("inventory paper %s:%s is not received by bank. Current state = %s", jeweler, paperNumber, paper.GetState())
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The repurchaser %q is ready to REPO the inventory financing paper  %q:%q. \nCurrent state = %q", paper.GetRepurchaser(), jeweler, paperNumber, paper.GetState())
	return paper, nil
}

//putInStorage updates a inventory paper to be put in storage
func (c *Contract) PutInStorage(ctx TransactionContextInterface, paperNumber string, jeweler string, supervisor string, storageAmount string, storageType string, storageAddress string, startdate string, endDate string, storageInfoUpdate string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)
	if err != nil {
		return nil, err
	}

	if paper.IsReceived() {
		if paper.GetSupervisor() == "" {
			paper.SetSupervisor(supervisor)
			paper.StorageAmount = storageAmount
			paper.StorageType = storageType
			paper.StorageAddress = storageAddress
			paper.StartDate = startdate
			paper.EndDate = endDate
			paper.StorageInfoUpdate = storageInfoUpdate

		}

	}

	if !paper.IsReceived() {
		return nil, fmt.Errorf("inventory paper %s:%s is not received by bank. Current state = %s", jeweler, paperNumber, paper.GetState())
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The supervisor %q has stored the inventory financing paper  %q:%q. \nThe info is storageAmount %q, storageType %q, storageAddress %q, endDate %q, storageInfoUpdate %q.The Current state = %q", paper.GetSupervisor(), jeweler, paperNumber, storageAmount, storageType, storageAddress, endDate, storageInfoUpdate, paper.GetState())
	return paper, nil
}

// Accept updates a inventory paper to be in accepted status and sets the next dealer
func (c *Contract) Accept(ctx TransactionContextInterface, paperNumber string, jeweler string, acceptedDateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.GetEvaluator() != "" && paper.GetRepurchaser() != "" && paper.GetSupervisor() != "" {
		paper.SetAccepted()
		paper.AcceptDateTime = acceptedDateTime

	}

	if !paper.IsAccepted() {
		return nil, fmt.Errorf("inventory paper %s:%s is not accepted by bank.The evaluator is %s. The repurchaser is %s. The supervisor is %s.Current state = %s", jeweler, paperNumber, paper.GetEvaluator(), paper.GetRepurchaser(), paper.GetSupervisor(), paper.GetState())
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The bank %q has accepted the inventory financing paper %q:%q .\nCurrent state is %q", paper.GetBank(), paper.GetEvaluator(), paperNumber, paper.GetState())
	return paper, nil
}

// Supervising updates a inventory paper to be in supervising status and sets the next dealer
func (c *Contract) Supervise(ctx TransactionContextInterface, paperNumber string, jeweler string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.IsAccepted() {
		paper.SetSupervising()
	}

	if !paper.IsSupervising() {
		return nil, fmt.Errorf("inventory paper %s:%s is not in supervision. Current state = %s", jeweler, paperNumber, paper.GetState())
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("inventory paper %q:%q is in supervision by %q. Current state = %q", jeweler, paperNumber, paper.GetSupervisor(), paper.GetState())
	return paper, nil
}

// Payback updates a inventory paper status to be paidback
func (c *Contract) Payback(ctx TransactionContextInterface, paperNumber string, jeweler string, paidBackDateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)
	if err != nil {
		return nil, err
	}

	if paper.IsPaidBack() {
		return nil, fmt.Errorf("paper %s:%s is already PaidBack", jeweler, paperNumber)
	}

	paper.SetPaidBack()
	paper.PaidbackDateTime = paidBackDateTime

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("inventory paper %q:%q is paid back by %q. Current state = %q", jeweler, paperNumber, jeweler, paper.GetState())
	return paper, nil
}

// Default updates a inventory paper status to be default
func (c *Contract) Default(ctx TransactionContextInterface, paperNumber string, jeweler string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.IsDefault() {
		return nil, fmt.Errorf("paper %s:%s has not been paidback", jeweler, paperNumber)
	}

	paper.SetDefault()

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("inventory paper %q:%q is not paid back by %q. Current state = %q", jeweler, paperNumber, jeweler, paper.GetState())
	return paper, nil
}

// Repurchase updates a inventory paper status to be repurchsed
func (c *Contract) Repurchase(ctx TransactionContextInterface, paperNumber string, jeweler string, repurchaseDateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.IsRepurchased() {
		return nil, fmt.Errorf("paper %s:%s is already repurchased", jeweler, paperNumber)
	}

	paper.SetRepurchased()
	paper.RepurchaseDateTime = repurchaseDateTime

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("inventory paper %q:%q is repurchased by %q. Current state = %q\n", jeweler, paperNumber, paper.GetRepurchaser(), paper.GetState())
	return paper, nil
}

// Reject a contract
func (c *Contract) Reject(ctx TransactionContextInterface, paperNumber string, jeweler string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if !paper.IsRejectable() {
		return nil, fmt.Errorf("paper %s:%s is not in rejectable state. CurrState: %s", jeweler, paperNumber, paper.GetState())
	}

	paper.LogPrevState()

	paper.SetApplied()

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}

	fmt.Printf("inventory paper %q:%q is rejected. Current state = %q\n", jeweler, paperNumber, paper.GetState())
	return paper, nil
}

// Revise a contract
func (c *Contract) Revise(ctx TransactionContextInterface, paperNumber string, jeweler string, financingAmount string, reviseDateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)
	if err != nil {
		return nil, err
	}

	if paper.GetState() != APPLIED {
		return nil, fmt.Errorf("paper %s:%s is not in applied state, CANNOT be revised. CurrState: %s", jeweler, paperNumber, paper.GetState())
	}

	paper.FinancingAmount = financingAmount
	paper.ReviseDateTime = reviseDateTime

	paper.Reinstate()

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The financing contract %s:%s is revised.\nCurrent Fin Amount is %q", jeweler, paperNumber, financingAmount)
	return paper, nil
}

// QueryPaper updates a inventory paper to be in received status and sets the next dealer
func (c *Contract) QueryPaper(ctx TransactionContextInterface, paperNumber string, jeweler string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}
	fmt.Printf("Current Paper: %q,%q.Current State = %q\n", jeweler, paperNumber, paper.GetState())
	return paper, nil
}

// // QueryAll returns  all papers found in world state
// func (c *Contract) QueryAllInventoryFinancingPapers(ctx TransactionContextInterface) ([]QueryResult, error) {
// 	list := ctx.GetPaperList()
// 	return list, nil

// 	// 	startKey := ""
// 	endKey := ""

// 	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	results := []QueryResult{}

// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()

// 		if err != nil {
// 			return nil, err
// 		}

// 		paper := new(InventoryFinancingPaper)
// 		_ = json.Unmarshal(queryResponse.Value, paper)

// 		queryResult := QueryResult{Key: queryResponse.Key, Record: paper}
// 		results = append(results, queryResult)
// 	}

// 	return results, nil
// }
