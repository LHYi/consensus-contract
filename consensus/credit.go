package consensus

import (
	"encoding/json"
	"fmt"

	ledgerapi "github.com/LHYi/consensus-contract/ledger-api"
)

// State enum for response credit state property
type State uint

const (
	// ISSUED state for when a credit has been issued
	ISSUED State = iota + 1
	// PENDING is the intermediate state of a trading process
	PENDING
	// TRADING state for when a credit is trading
	TRADING
	// REDEEMED state for when a credit has been redeemed
	REDEEMED
)

func (state State) String() string {
	names := []string{"ISSUED", "PENDING", "TRADING", "REDEEMED"}

	if state < ISSUED || state > REDEEMED {
		return "UNKNOWN"
	}

	return names[state-1]
}

// CreateResponseCreditKey creates a key for response credits
func CreateResponseCreditKey(issuer string, creditNumber string) string {
	return ledgerapi.MakeKey(issuer, creditNumber)
}

// Used for managing the fact status is private but want it in world state
type responseCreditAlias ResponseCredit
type jsonResponseCredit struct {
	*responseCreditAlias
	State State  `json:"currentState"`
	Class string `json:"class"`
	Key   string `json:"key"`
}

// TODO: change the properties accordingly
// ResponseCredit defines a response credit
type ResponseCredit struct {
	CreditNumber  string `json:"creditNumber"`
	Issuer        string `json:"issuer"`
	IssuerMSP     string `json:"issuerMSP"`
	IssueDateTime string `json:"issueDateTime"`
	//FaceValue        int    `json:"faceValue"`
	//MaturityDateTime string `json:"maturityDateTime"`
	Owner    string `json:"owner"`
	OwnerMSP string `json:"OwnerMSP"`
	state    State  `metadata:"currentState"`
	class    string `metadata:"class"`
	key      string `metadata:"key"`
}

// UnmarshalJSON special handler for managing JSON marshalling
func (rc *ResponseCredit) UnmarshalJSON(data []byte) error {
	jrc := jsonResponseCredit{responseCreditAlias: (*responseCreditAlias)(rc)}

	err := json.Unmarshal(data, &jrc)

	if err != nil {
		return err
	}

	rc.state = jrc.State

	return nil
}

// MarshalJSON special handler for managing JSON marshalling
func (rc ResponseCredit) MarshalJSON() ([]byte, error) {
	jrc := jsonResponseCredit{responseCreditAlias: (*responseCreditAlias)(&rc), State: rc.state, Class: "org.creditnet.responsecredit", Key: ledgerapi.MakeKey(rc.Issuer, rc.CreditNumber)}

	return json.Marshal(&jrc)
}

// GetState returns the state
func (rc *ResponseCredit) GetState() State {
	return rc.state
}

// SetIssued returns the state to issued
func (rc *ResponseCredit) SetIssued() {
	rc.state = ISSUED
}

// SetPending sets the state to pending
func (rc *ResponseCredit) SetPending() {
	rc.state = PENDING
}

// SetTrading sets the state to trading
func (rc *ResponseCredit) SetTrading() {
	rc.state = TRADING
}

// SetRedeemed sets the state to redeemed
func (rc *ResponseCredit) SetRedeemed() {
	rc.state = REDEEMED
}

// IsIssued returns true if state is issued
func (rc *ResponseCredit) IsIssued() bool {
	return rc.state == ISSUED
}

// IsPending returns true if state is pending
func (rc *ResponseCredit) IsPending() bool {
	return rc.state == PENDING
}

// IsTrading returns true if state is trading
func (rc *ResponseCredit) IsTrading() bool {
	return rc.state == TRADING
}

// IsRedeemed returns true if state is redeemed
func (rc *ResponseCredit) IsRedeemed() bool {
	return rc.state == REDEEMED
}

// GetSplitKey returns values which should be used to form key
func (rc *ResponseCredit) GetSplitKey() []string {
	return []string{rc.Issuer, rc.CreditNumber}
}

// Serialize formats the commercial paper as JSON bytes
func (rc *ResponseCredit) Serialize() ([]byte, error) {
	return json.Marshal(rc)
}

// Deserialize formats the commercial paper from JSON bytes
func Deserialize(bytes []byte, rc *ResponseCredit) error {
	err := json.Unmarshal(bytes, rc)

	if err != nil {
		return fmt.Errorf("Error deserializing response credit. %s", err.Error())
	}

	return nil
}
