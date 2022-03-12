package consensus

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionContextInterface an interface to
// describe the minimum required functions for
// a transaction context in the response
// credit
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetCreditList() ListInterface
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// response credit contract
type TransactionContext struct {
	contractapi.TransactionContext
	creditList *list
}

// GetCreditList return credit list
func (tc *TransactionContext) GetCreditList() ListInterface {
	if tc.creditList == nil {
		tc.creditList = newList(tc)
	}

	return tc.creditList
}
