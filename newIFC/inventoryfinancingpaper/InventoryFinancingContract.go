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
			Jeweler:               "Prius",
			ApplyDateTime:         "202109",
			FinancingAmount:       "100000",
			Productor:             "kaka",
			ProductType:           "Gold",
			ProductAmount:         "4000000",
			ProductDate:           "202106",
			ProductInfoUpdateTime: "202111",
			BrandCompany:          "Chou Tai Fork",
			GrantedObject:         "cici retailer",
			GrantedStartDate:      "202001",
			GrantedEndDate:        "202203",
			GrantedInfoUpdateTime: "202111",
			AuthorizedDate:        "202110",
			Bank:                  "Heng Seng",
			ReceiveDateTime:       "202111",
			Evaluator:             "coco",
			EvalDateTime:          "202112",
			EvalType:              "Gold",
			EvalQualityProportion: "99%",
			EvalAmount:            "10000",
			Supervisor:            "Pika",
			StorageAmount:         "40000",
			StorageType:           "Gold",
			StorageAddress:        "Earth street ka",
			EndDate:               "202203",
			StorageInfoUpdate:     "202111",
			Repurchaser:           "Chou Seng Seng",
			AcceptedDateTime:      "202111",
			PaidbackDateTime:      "",
			RepurchaseDateTime:    "202203",
			state:                 PAIDBACK,
			prevstate:             SUPERVISING,
			class:                 "",
			key:                   "",
		},
		InventoryFinancingPaper{
			PaperNumber:           "222",
			Jeweler:               "Mustang",
			ApplyDateTime:         "202110",
			FinancingAmount:       "222222",
			Productor:             "haha",
			ProductType:           "Diamond",
			ProductAmount:         "20000",
			ProductDate:           "202006",
			ProductInfoUpdateTime: "202101",
			BrandCompany:          "Chou Seng Seng",
			GrantedObject:         "retailer 2",
			GrantedStartDate:      "202102",
			GrantedEndDate:        "202109",
			GrantedInfoUpdateTime: "202111",
			AuthorizedDate:        "202109",
			Bank:                  "China Bank",
			ReceiveDateTime:       "202111",
			Evaluator:             "kimi",
			EvalDateTime:          "202111",
			EvalType:              "Diamond",
			EvalQualityProportion: "98%",
			EvalAmount:            "250000",
			Supervisor:            "pika",
			StorageAmount:         "250000",
			StorageType:           "Diamond",
			StorageAddress:        "Jupiter",
			EndDate:               "202211",
			StorageInfoUpdate:     "202111",
			Repurchaser:           "Chou Tai Fork",
			AcceptedDateTime:      "202211",
			PaidbackDateTime:      "202201",
			RepurchaseDateTime:    "",
			state:                 REPURCHADED,
			prevstate:             SUPERVISING,
			class:                 "",
			key:                   "",
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
func (c *Contract) Apply(ctx TransactionContextInterface, paperNumber string, jeweler string, financingAmount string, applyDateTime string) (*InventoryFinancingPaper, error) {
	paper := InventoryFinancingPaper{PaperNumber: paperNumber, Jeweler: jeweler, FinancingAmount: financingAmount, ApplyDateTime: applyDateTime}

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
func (c *Contract) OfferLisenceInfo(ctx TransactionContextInterface, paperNumber string, jeweler string, brandCompany string, grantedObject string, grantedStartDate string, grantedEndDate string, grantedInfoUpdateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)
	if err != nil {
		return nil, err
	}
	if paper.IsApplied() {
		if paper.GetBrandCompany() == "" {
			paper.SetBrandCompany(brandCompany)
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
func (c *Contract) Evaluate(ctx TransactionContextInterface, paperNumber string, jeweler string, evaluator string, evalType string, evalQualityProportion string, evalAmount string, evalDateTime string) (*InventoryFinancingPaper, error) {
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
func (c *Contract) PutInStorage(ctx TransactionContextInterface, paperNumber string, jeweler string, supervisor string, storageAmount string, storageType string, storageAddress string, endDate string, storageInfoUpdate string) (*InventoryFinancingPaper, error) {
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
		paper.AcceptedDateTime = acceptedDateTime

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
func (c *Contract) Supervise(ctx TransactionContextInterface, paperNumber string, jeweler string, supervisor string) (*InventoryFinancingPaper, error) {
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
func (c *Contract) Repurchase(ctx TransactionContextInterface, jeweler string, paperNumber string, repurchaseDateTime string) (*InventoryFinancingPaper, error) {
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
func (c *Contract) Reject(ctx TransactionContextInterface, jeweler string, paperNumber string) (*InventoryFinancingPaper, error) {
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
func (c *Contract) Revise(ctx TransactionContextInterface, jeweler string, paperNumber string, financingAmount string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)
	if err != nil {
		return nil, err
	}

	if paper.GetState() != APPLIED {
		return nil, fmt.Errorf("paper %s:%s is not in applied state, CANNOT be revised. CurrState: %s", jeweler, paperNumber, paper.GetState())
	}

	paper.FinancingAmount = financingAmount

	paper.Reinstate()

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The financing contract %s:%s is revised.\nCurrent Fin Amount is %q", jeweler, paperNumber, financingAmount)
	return paper, nil
}

// QueryPaper updates a inventory paper to be in received status and sets the next dealer
func (c *Contract) QueryPaper(ctx TransactionContextInterface, jeweler string, paperNumber string) (*InventoryFinancingPaper, error) {
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
