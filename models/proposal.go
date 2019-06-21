package models

import tezos "github.com/ecadlabs/go-tezos"

type ProposalSummary struct {
	tezos.Proposal
	Cycle         int
	NewSupporters int
}

type Proposal struct {
	ProposalHash string
	PKH          string
	Period       int
}
