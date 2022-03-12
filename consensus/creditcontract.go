package consensus

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract chaincode that defines
// the business logic for managing response credit

type Contract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (c *Contract) Instantiate() {
	fmt.Println("Instantiated")
}

// Issue creates a new response credit and stores it in the world state
func (c *Contract) Issue(ctx TransactionContextInterface, creditNumber string, issuer string, issueDateTime string) (*ResponseCredit, error) {
	clientIdentity, err0 := ctx.GetClientIdentity().GetMSPID()
	if err0 != nil {
		return nil, fmt.Errorf("Failed to get client identity")
	}

	credit := ResponseCredit{CreditNumber: creditNumber, Issuer: issuer, IssuerMSP: clientIdentity, IssueDateTime: issueDateTime, Owner: issuer, OwnerMSP: clientIdentity}
	credit.SetIssued()

	err := ctx.GetCreditList().AddCredit(&credit)

	if err != nil {
		return nil, err
	}

	return &credit, nil
}

// Query returns the credit queried by the given issuer and credir number
func (c *Contract) Query(ctx TransactionContextInterface, creditNumber string, issuer string) (*ResponseCredit, error) {
	credit, err := ctx.GetCreditList().GetCredit(issuer, creditNumber)

	if err != nil {
		return nil, err
	}

	return credit, nil
}

// Buy updates a response credit to be in trading status and sets the new owner
func (c *Contract) Buy(ctx TransactionContextInterface, issuer string, creditNumber string, currentOwner string, currentOwnerMSP string, newOwner string, newOwnerMSP string, price int, purchaseDateTime string) (*ResponseCredit, error) {
	credit, err := ctx.GetCreditList().GetCredit(issuer, creditNumber)

	if err != nil {
		return nil, err
	}

	if credit.Owner != currentOwner {
		return nil, fmt.Errorf("Credit %s:%s is not owned by %s", issuer, creditNumber, currentOwner)
	}

	if credit.IsIssued() {
		credit.SetTrading()
	}

	if !credit.IsTrading() {
		return nil, fmt.Errorf("Credit %s:%s is not trading. Current state = %s", issuer, creditNumber, credit.GetState())
	}

	credit.Owner = newOwner

	err = ctx.GetCreditList().UpdateCredit(credit)

	if err != nil {
		return nil, err
	}

	return credit, nil
}

// BuyRequest is invoked by the buyer, which checks the owner of the response credit and set it to the pending state
func (c *Contract) BuyRequest(ctx TransactionContextInterface, issuer string, creditNumber string, currentOwner string, newOwner string, price int, purchaseDateTime string) (*ResponseCredit, error) {
	credit, err := ctx.GetCreditList().GetCredit(issuer, creditNumber)

	if err != nil {
		return nil, err
	}

	if credit.IsIssued() {
		credit.SetTrading()
	}

	if !credit.IsTrading() {
		return nil, fmt.Errorf("Credit %s:%s is not trading. Current state = %s", issuer, creditNumber, credit.GetState())
	}

	credit.SetPending()

	err = ctx.GetCreditList().UpdateCredit(credit)

	if err != nil {
		return nil, err
	}

	return credit, nil
}

// Transfer function only allows the owner of the commercial paper to execute, which is the complement to the BuyRequest function.
func (c *Contract) Transfer(ctx TransactionContextInterface, issuer string, creditNumber string, currentOwner string, newOwner string, newOwnerMSP string, confirmDateTime string) (*ResponseCredit, error) {
	credit, err := ctx.GetCreditList().GetCredit(issuer, creditNumber)

	clientIdentity, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("Failed to get client identity")
	}
	if credit.OwnerMSP != clientIdentity {
		return nil, fmt.Errorf("Credit %s:%s is not owned by you, failed to transfer", issuer, creditNumber)
	}

	if !credit.IsPending() {
		return nil, fmt.Errorf("Credit %s:%s is not pending. Current state = %s. Need to invoke BuyRequest first.", issuer, creditNumber, credit.GetState())
	}

	credit.Owner = newOwner
	credit.OwnerMSP = newOwnerMSP
	credit.SetTrading()

	err = ctx.GetCreditList().UpdateCredit(credit)

	if err != nil {
		return nil, err
	}

	return credit, nil
}

// Redeem updates a response credit status to be redeemed
func (c *Contract) Redeem(ctx TransactionContextInterface, issuer string, creditNumber string, redeemingOwner string, redeenDateTime string) (*ResponseCredit, error) {
	credit, err := ctx.GetCreditList().GetCredit(issuer, creditNumber)

	if err != nil {
		return nil, err
	}

	if credit.IsRedeemed() {
		return nil, fmt.Errorf("Credit %s:%s is already redeemed", issuer, creditNumber)
	}

	clientIdentity, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("Failed to get client identity")
	}
	if credit.OwnerMSP != clientIdentity {
		return nil, fmt.Errorf("Credit %s:%s is not owned by you, failed to invoke Redeem", issuer, creditNumber)
	}

	if credit.Owner != redeemingOwner {
		return nil, fmt.Errorf("Credit %s:%s is not owned by %s", issuer, creditNumber, redeemingOwner)
	}

	credit.Owner = credit.Issuer
	credit.OwnerMSP = credit.IssuerMSP
	credit.SetRedeemed()

	err = ctx.GetCreditList().UpdateCredit(credit)

	if err != nil {
		return nil, err
	}

	return credit, nil
}
