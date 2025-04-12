package message

type CosmosAPIResponse struct {
	DelegationResponses []DelegationResponse `json:"delegation_responses"`
}

type DelegationResponse struct {
	Delegation Delegation `json:"delegation"`
	Balance    Balance    `json:"balance"`
}

type Delegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
}

type Balance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
