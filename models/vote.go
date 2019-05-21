package models

// Ballot is a struct holding tezos ballot information
type Ballot struct {
	PKH          string
	Ballot       string
	ProposalHash string
	Rolls        int64
}
