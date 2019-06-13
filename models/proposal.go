package models

import tezos "github.com/ecadlabs/go-tezos"

type ProposalSummary = tezos.Proposal

type Proposal struct {
	ProposalHash string
	PKH          string
	Period       int
}
