package consensus

import ledgerapi "github.com/LHYi/consensus-contract/ledger-api"

// ListInterface defines functionality needed
// to interact with the world state on behalf
// of a response credit
type ListInterface interface {
	AddCredit(*ResponseCredit) error
	GetCredit(string, string) (*ResponseCredit, error)
	UpdateCredit(*ResponseCredit) error
}

type list struct {
	stateList ledgerapi.StateListInterface
}

func (rcl *list) AddCredit(credit *ResponseCredit) error {
	return rcl.stateList.AddState(credit)
}

func (rcl *list) GetCredit(issuer string, creditNumber string) (*ResponseCredit, error) {
	rc := new(ResponseCredit)

	err := rcl.stateList.GetState(CreateResponseCreditKey(issuer, creditNumber), rc)

	if err != nil {
		return nil, err
	}

	return rc, nil
}

func (rcl *list) UpdateCredit(credit *ResponseCredit) error {
	return rcl.stateList.UpdateState(credit)
}

// NewList create a new list from context
func newList(ctx TransactionContextInterface) *list {
	stateList := new(ledgerapi.StateList)
	stateList.Ctx = ctx
	stateList.Name = "org.creditnet.responsecreditlist"
	stateList.Deserialize = func(bytes []byte, state ledgerapi.StateInterface) error {
		return Deserialize(bytes, state.(*ResponseCredit))
	}

	list := new(list)
	list.stateList = stateList

	return list
}
