/*
 * SPDX-License-Identifier: Apache-2.0
 */

package inventoryfinancingpaper

import (
	"encoding/json"
	"fmt"

	ledgerapi "github.com/hyperledger/fabric-samples/commercial-paper/organization/digibank/contract-go/ledger-api"
)

// State enum for inventory financing  state property
//并行功能不设置状态
type State uint

const (
	APPLIED = iota + 1
	RECEIVED
	ACCEPTED
	SUPERVISING
	PAIDBACK
	DEFAULT
	REPURCHADED
)

func (state State) String() string {
	names := []string{"APPLIED ", "RECEIVED", "ACCEPTED", "SUPERVISING", "PAIDBACK", "DEFAULT", "REPURCHADED"}

	if state < APPLIED || state > REPURCHADED {
		return "UNKNOWN"
	}

	return names[state-1]
}

// CreateInventoryFinancingPaperKey creates a key for inventory financing
func CreateInventoryFinancingPaperKey(jeweler string, paperNumber string) string {
	return ledgerapi.MakeKey(jeweler, paperNumber)
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *InventoryFinancingPaper
}

// Used for managing the fact status is private but want it in world state
type InventoryFinancingPaperAlias InventoryFinancingPaper
type jsonInventoryFinancingPaper struct {
	*InventoryFinancingPaperAlias
	State State  `json:"currentState"`
	Class string `json:"class"`
	Key   string `json:"key"`
}

// InventoryFinancingPaper 定义了一个珠宝存货融资流程
type InventoryFinancingPaper struct {
	//珠宝商发起融资申请
	PaperNumber     string `json:"paperNumber"`     //融资申请编号
	Jeweler         string `json:"jeweler"`         //融资申请珠宝商
	ApplyDateTime   string `json:"applyDateTime"`   //提交申请时间（web端自动生成）
	FinancingAmount string `json:"financingAmount"` //融资金额
	//生产者提供了生产信息上链
	Productor             string `json:"productor"`             //生产商
	ProductType           string `json:"productType"`           //货品种类
	ProductAmount         string `json:"productAmount"`         //货品数量
	ProductDate           string `json:"productDate"`           //货品生产日期
	ProductInfoUpdateTime string `json:"productInfoUpdateTime"` //货品信息更新日期（web端自动生成）
	//品牌方提供授权信息上链
	BrandCompany          string `json:"brandCompany "`
	GrantedObject         string `json:"grantedObject "`        //授权对象
	GrantedStartDate      string `json:"grantedStartDate"`      //授权开始日期
	GrantedEndDate        string `json:"grantedEndDate"`        //授权结束日期
	GrantedInfoUpdateTime string `json:"grantedInfoUpdateTime"` //授权信息更新日期（web端自动生成）
	//银行认证供应链各方的背书
	AuthorizedDate string `json:"authorizedDate "` //认证和授权时间（web端自动生成）
	//银行收到融资申请
	Bank            string `json:"bank"`
	ReceiveDateTime string `json:"receiveDateTime"` //收到融资申请时间（web端自动生成）
	//评估鉴定方提供鉴定信息
	Evaluator             string `json:"evaluator"`
	EvalDateTime          string `json:"evalDateTime"`          //鉴定时间（web端自动生成）
	EvalType              string `json:"evalType"`              //评估种类
	EvalQualityProportion string `json:"evalQualityProportion"` //评估质量（质检合格比例）
	EvalAmount            string `json:"evalAmount"`            //评估价值
	//仓库监管方提供仓单信息
	Supervisor        string `json:"supervisor"`
	StorageAmount     string `json:"storageAmount"`     //仓库货品总量
	StorageType       string `json:"storageType"`       //货品种类
	StorageAddress    string `json:"storageAddress"`    //存储地址
	EndDate           string `json:"endDate"`           //期限
	StorageInfoUpdate string `json:"storageInfoUpdate"` //出具仓单的时间（web端自动生成）
	//回购方准备好可以后续回购
	Repurchaser   string `json:"repurchaser"`
	ReadyDateTime string `json:"readyDateTime"`
	//银行接受
	AcceptedDateTime string `json:"acceptedDateTime"` //银行接受时间（web端自动生成）
	//珠宝商回购
	PaidbackDateTime string `json:"paidBackDateTime"` //珠宝商回购时间（web端自动生成）
	//回购方回购
	RepurchaseDateTime string `json:"repurchaseDateTime"` //回购方回购时间（web端自动生成）
	state              State  `metadata:"currentState"`
	prevstate          State  `metadata:"prevState,optional"`
	class              string `metadata:"class"`
	key                string `metadata:"key"`
	//珠宝商重写
	ReviseDateTime string `json:"reviseDateTime"`
}

// UnmarshalJSON special handler for managing JSON marshalling
func (ifc *InventoryFinancingPaper) UnmarshalJSON(data []byte) error {
	jifc := jsonInventoryFinancingPaper{InventoryFinancingPaperAlias: (*InventoryFinancingPaperAlias)(ifc)}

	err := json.Unmarshal(data, &jifc)

	if err != nil {
		return err
	}

	ifc.state = jifc.State

	return nil
}

// MarshalJSON special handler for managing JSON marshalling
func (ifc InventoryFinancingPaper) MarshalJSON() ([]byte, error) {
	jifc := jsonInventoryFinancingPaper{InventoryFinancingPaperAlias: (*InventoryFinancingPaperAlias)(&ifc), State: ifc.state, Class: "org.papernet.InventoryFinancingPaper", Key: ledgerapi.MakeKey(ifc.Jeweler, ifc.PaperNumber)}

	return json.Marshal(&jifc)
}

// GetState returns the state
func (ifc *InventoryFinancingPaper) GetState() State {
	return ifc.state
}

// SetPrevState returns the previous state
func (ifc *InventoryFinancingPaper) LogPrevState() State {
	ifc.prevstate = ifc.state
	return ifc.prevstate
}

// Get prev state and set as curr state
func (ifc *InventoryFinancingPaper) Reinstate() State {
	ifc.state = ifc.prevstate
	return ifc.state
}

// GetBank returns the bank
func (ifc *InventoryFinancingPaper) GetBank() string {
	return ifc.Bank
}

// SetBank set the Bank to bank
func (ifc *InventoryFinancingPaper) SetBank(bank string) {
	ifc.Bank = bank
}

// GetProductor returns the productor
func (ifc *InventoryFinancingPaper) GetProductor() string {
	return ifc.Productor
}

// SetProductor set the Productor to productor
func (ifc *InventoryFinancingPaper) SetProductor(productor string) {
	ifc.Productor = productor
}

// GetBrandCompany returns the productor
func (ifc *InventoryFinancingPaper) GetBrandCompany() string {
	return ifc.BrandCompany
}

// SetBrandCompany set the BrandCompany to brandCompany
func (ifc *InventoryFinancingPaper) SetBrandCompany(brandCompany string) {
	ifc.BrandCompany = brandCompany
}

// GetEvaluator returns the evaluator
func (ifc *InventoryFinancingPaper) GetEvaluator() string {
	return ifc.Evaluator
}

// SetEvaluator set the Evaluator to evaluator
func (ifc *InventoryFinancingPaper) SetEvaluator(evaluator string) {
	ifc.Evaluator = evaluator
}

// GetRepurchaser returns the repurchaser
func (ifc *InventoryFinancingPaper) GetRepurchaser() string {
	return ifc.Repurchaser
}

// SetRepurchaser set the Repurchaser to repurchaser
func (ifc *InventoryFinancingPaper) SetRepurchaser(repurchaser string) {
	ifc.Repurchaser = repurchaser
}

// GetSupervisor returns the supervisor
func (ifc *InventoryFinancingPaper) GetSupervisor() string {
	return ifc.Supervisor
}

// SetSupervisor set the state to supervisor
func (ifc *InventoryFinancingPaper) SetSupervisor(supervisor string) {
	ifc.Supervisor = supervisor
}

// GetEndDate returns the receivedatetime
func (ifc *InventoryFinancingPaper) GetEndDate() string {
	return ifc.EndDate
}

// SetEndDate set the EndDate to endDate
func (ifc *InventoryFinancingPaper) SetEndDate(endDate string) {
	ifc.EndDate = endDate
}

// SetApplied set the state to applied
func (ifc *InventoryFinancingPaper) SetApplied() {
	ifc.state = APPLIED
}

// SetReceived sets the state to received
func (ifc *InventoryFinancingPaper) SetReceived() {
	ifc.state = RECEIVED
}

// SetAccepted sets the state to accepted
func (ifc *InventoryFinancingPaper) SetAccepted() {
	ifc.state = ACCEPTED
}

// SetSupervising sets the state to supervising
func (ifc *InventoryFinancingPaper) SetSupervising() {
	ifc.state = SUPERVISING
}

// SetPaidBack sets the state to paidBack
func (ifc *InventoryFinancingPaper) SetPaidBack() {
	ifc.state = PAIDBACK
}

// SetDefault sets the state to default
func (ifc *InventoryFinancingPaper) SetDefault() {
	ifc.state = DEFAULT
}

// SetDefault sets the state to repurchased
func (ifc *InventoryFinancingPaper) SetRepurchased() {
	ifc.state = REPURCHADED
}

// IsApplied returns true if state is applied
func (ifc *InventoryFinancingPaper) IsApplied() bool {
	return ifc.state == APPLIED
}

// IsReceived returns true if state is received
func (ifc *InventoryFinancingPaper) IsReceived() bool {
	return ifc.state == RECEIVED
}

// IsAccepted returns true if state is accepted
func (ifc *InventoryFinancingPaper) IsAccepted() bool {
	return ifc.state == ACCEPTED
}

// Supervising returns true if state is supervising
func (ifc *InventoryFinancingPaper) IsSupervising() bool {
	return ifc.state == SUPERVISING
}

// IsPaidBack returns true if state is paidback
func (ifc *InventoryFinancingPaper) IsPaidBack() bool {
	return ifc.state == PAIDBACK
}

// IsDefault returns true if state is default
func (ifc *InventoryFinancingPaper) IsDefault() bool {
	return ifc.state == DEFAULT
}

// IsRepurchased returns true if state is repurchased
func (ifc *InventoryFinancingPaper) IsRepurchased() bool {
	return ifc.state == REPURCHADED
}

// IsRejectable returns true if state is in RECEIVED
func (ifc *InventoryFinancingPaper) IsRejectable() bool {
	var ret bool = false

	if ifc.state == RECEIVED {
		ret = true
	}
	return ret
}

// GetSplitKey returns values which should be used to form key
func (ifc *InventoryFinancingPaper) GetSplitKey() []string {
	return []string{ifc.Jeweler, ifc.PaperNumber}
}

// Serialize formats the inventory financing  as JSON bytes
func (ifc *InventoryFinancingPaper) Serialize() ([]byte, error) {
	return json.Marshal(ifc)
}

// Deserialize formats the Inventory Financing  from JSON bytes
func Deserialize(bytes []byte, ifc *InventoryFinancingPaper) error {
	err := json.Unmarshal(bytes, ifc)

	if err != nil {
		return fmt.Errorf("error deserializing inventory financing . %s", err.Error())
	}

	return nil
}
