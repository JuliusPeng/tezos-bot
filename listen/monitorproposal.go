package listen

import (
	"context"
	"fmt"
	"log"

	tezos "github.com/ecadlabs/go-tezos"
	"github.com/ecadlabs/tezos-bot/models"
)

func (t *TezosListener) lookForProposal(ctx context.Context, block *tezos.Block) error {
	log.Printf("TezosListener: Inspecting block %s for new proposal operations.\n", block.Hash)

	// Retrieve proposals for the current phase
	existingProposals, err := t.service.GetProposals(ctx, t.config.GetChainID(), fmt.Sprintf("%d", block.Header.Level-1))

	if err != nil {
		return err
	}

	proposalExists := func(proposal string) bool {
		for _, existingProposal := range existingProposals {
			if existingProposal.ProposalHash == proposal {
				return true
			}
		}
		return false
	}

	if err != nil {
		return err
	}

	proposalOps := []*tezos.ProposalOperationElem{}
	for _, group := range block.Operations {
		for _, op := range group {
			proposalOps = append(proposalOps, tezos.FilterProposalOps(op.Contents)...)
		}
	}

	listings, err := t.service.GetBallotListings(ctx, t.config.GetChainID(), block.Hash)

	if err != nil {
		return err
	}

	totalRolls := int64(0)
	for _, entry := range listings {
		totalRolls += entry.Rolls
	}

	if totalRolls == 0 {
		// Unlikely to occurs
		return fmt.Errorf("No rolls found in this block")
	}

	for _, proposalOp := range proposalOps {
		for _, proposal := range proposalOp.Proposals {
			rolls := int64(0)
			for _, entry := range listings {
				if entry.PKH == proposalOp.Source {
					rolls = entry.Rolls
				}
			}
			p := &models.Proposal{
				ProposalHash: proposal,
				PKH:          proposalOp.Source,
				Period:       proposalOp.Period,
				Rolls:        rolls,
			}

			if !proposalExists(proposal) {
				t.newProposalChan <- p
			} else {
				t.proposalUpvoteChan <- p
			}
		}
	}

	return nil
}
