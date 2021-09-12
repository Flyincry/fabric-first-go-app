/*
 * SPDX-License-Identifier: Apache-2.0
 */

 package inventoryfinancingpaper

import ledgerapi "github.com/hyperledger/fabric-samples/commercial-paper/organization/digibank/contract-go/ledger-api"

// ListInterface defines functionality needed
// to interact with the world state on behalf
// of a commercial paper
type ListInterface interface {
	AddPaper(*InventoryFinancingPaper) error
	GetPaper(string, string) (*InventoryFinancingPaper, error)
	UpdatePaper(*InventoryFinancingPaper) error
}

type list struct {
	stateList ledgerapi.StateListInterface
}

func (ifcl *list) AddPaper(paper *InventoryFinancingPaper) error {
	return ifcl.stateList.AddState(paper)
}

func (ifcl *list) GetPaper(jeweler string, paperNumber string) (*InventoryFinancingPaper, error) {
	ifc := new(InventoryFinancingPaper)

	err := ifcl.stateList.GetState(CreateInventoryFinancingPaperKey(jeweler, paperNumber), ifc)

	if err != nil {
		return nil, err
	}

	return ifc, nil
}

func (ifcl *list) UpdatePaper(paper *InventoryFinancingPaper) error {
	return ifcl.stateList.UpdateState(paper)
}

// NewList create a new list from context
func newList(ctx TransactionContextInterface) *list {
	stateList := new(ledgerapi.StateList)
	stateList.Ctx = ctx
	stateList.Name = "org.inventoryfinancingpaperlist"
	stateList.Deserialize = func(bytes []byte, state ledgerapi.StateInterface) error {
		return Deserialize(bytes, state.(*InventoryFinancingPaper))
	}

	list := new(list)
	list.stateList = stateList

	return list
}