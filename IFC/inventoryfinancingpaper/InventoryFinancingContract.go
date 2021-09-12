/*
 * SPDX-License-Identifier: Apache-2.0
 */

package inventoryfinancingpaper

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract chaincode that defines
// the business logic for managing inventory paper

type Contract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (c *Contract) Init() {
	fmt.Println("Instantiated")
}

// Apply creates a new inventory paper and stores it in the world state
func (c *Contract) Apply(ctx TransactionContextInterface, paperNumber string, jeweler string, applyDateTime string, financingAmount int) (*InventoryFinancingPaper, error) {
	paper := InventoryFinancingPaper{PaperNumber: paperNumber, Jeweler: jeweler, ApplyDateTime: applyDateTime, FinancingAmount: financingAmount}
	paper.SetApplied()
	paper.LogPrevState()

	err := ctx.GetPaperList().AddPaper(&paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The jeweler %q  has applied for a new inventory financingp paper %q, the apply date is %q,the financing amount is %v.\n Current State is %q", jeweler, paperNumber, applyDateTime, financingAmount, paper.GetState())
	return &paper, nil
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

// Receive updates a inventory paper to be in received status and sets the next dealer
func (c *Contract) Receive(ctx TransactionContextInterface, jeweler string, bank string, paperNumber string, receiveDateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.GetBank() == "" {
		paper.SetBank(bank)
	}

	if paper.IsApplied() {
		paper.SetReceived()
	}

	if !paper.IsReceived() {
		return nil, fmt.Errorf("inventory paper %s:%s is not received by bank. Current state = %s", jeweler, paperNumber, paper.GetState())
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The bank %q has received the inventory financing paper %q from jeweler %q,the receive date is %q \n Current State is %q", paper.GetBank(), paperNumber, jeweler, receiveDateTime, paper.GetState())
	return paper, nil
}

//Evaluated updates a inventory paper to be in Evaluated status and sets the next dealer
func (c *Contract) Evaluate(ctx TransactionContextInterface, jeweler string, paperNumber string, evaluator string, evalDateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.GetEvaluator() == "" {
		paper.SetEvaluator(evaluator)
	}

	if paper.IsReceived() {
		paper.SetEvaluated()
	}

	if !paper.IsEvaluated() {
		return nil, fmt.Errorf("inventory paper %s:%s is not yet evaluated, Current state = %s", jeweler, paperNumber, paper.GetState())
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The evluator %q has evaluated the inventory financing paper %q:%q,the evaluate date is %q..\n Current State is %q", paper.GetEvaluator(), jeweler, paperNumber, evalDateTime, paper.GetState())
	return paper, nil
}

//ReadyRepo updates a inventory paper to be in ReadyRepo status and sets the next dealer
func (c *Contract) ReadyRepo(ctx TransactionContextInterface, jeweler string, paperNumber string, repurchaser string, readyDateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.GetRepurchaser() == "" {
		paper.SetRepurchaser(repurchaser)
	}

	if paper.IsEvaluated() {
		paper.SetReadyREPO()
	}

	if !paper.IsReadyREPO() {
		return nil, fmt.Errorf("inventory paper %q:%q is waiting for REPO's ready. Current state = %q", jeweler, paperNumber, paper.GetState())
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The repurchaser %q is ready to REPO the inventory financing paper  %q:%q, the ready date is %q.\nCurrent state = %q", paper.GetRepurchaser(), jeweler, paperNumber, readyDateTime, paper.GetState())
	return paper, nil
}

// Accept updates a inventory paper to be in accepted status and sets the next dealer
func (c *Contract) Accept(ctx TransactionContextInterface, jeweler string, paperNumber string, acceptDate string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.IsReadyREPO() {
		paper.SetAccepted()
	}

	if !paper.IsAccepted() {
		return nil, fmt.Errorf("inventory paper %s:%s is not accepted by bank. Current state = %s", jeweler, paperNumber, paper.GetState())
	}

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("The bank %q has accepted the inventory financing paper %q:%q ,The accept date is %q.\nCurrent state is %q", paper.GetBank(), paper.GetEvaluator(), paperNumber, acceptDate, paper.GetState())
	return paper, nil
}

// Supervising updates a inventory paper to be in supervising status and sets the next dealer
func (c *Contract) Supervise(ctx TransactionContextInterface, jeweler string, supervisor string, endDate string, paperNumber string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.GetSupervisor() == "" {
		paper.SetSupervisor(supervisor)
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
	fmt.Printf("inventory paper %q:%q is in supervision by %q,The end date is %q. Current state = %q", jeweler, paperNumber, paper.GetSupervisor(), endDate, paper.GetState())
	return paper, nil
}

// Payback updates a inventory paper status to be paidback
func (c *Contract) Payback(ctx TransactionContextInterface, jeweler string, paperNumber string, paidbackDateTime string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.IsPaidBack() {
		return nil, fmt.Errorf("paper %s:%s is already PaidBack", jeweler, paperNumber)
	}

	paper.SetPaidBack()

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("inventory paper %q:%q is paid back by %q,The paidback date is %q. Current state = %q", jeweler, paperNumber, jeweler, paidbackDateTime, paper.GetState())
	return paper, nil
}

// Default updates a inventory paper status to be default
func (c *Contract) Default(ctx TransactionContextInterface, jeweler string, paperNumber string) (*InventoryFinancingPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(jeweler, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.IsDefault() {
		return nil, fmt.Errorf("paper %s:%s can not be paidback", jeweler, paperNumber)
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
		return nil, fmt.Errorf("paper %s:%s is already Repurchased", jeweler, paperNumber)
	}

	paper.SetRepurchased()

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	fmt.Printf("inventory paper %q:%q is repurchased by %q,The repurchased date is %q. Current state = %q\n", jeweler, paperNumber, paper.GetRepurchaser(), repurchaseDateTime, paper.GetState())
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
func (c *Contract) Revise(ctx TransactionContextInterface, jeweler string, paperNumber string, reviseDateTime string, financingAmount int) (*InventoryFinancingPaper, error) {
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
