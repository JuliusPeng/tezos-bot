package listen

import (
	"context"
	"fmt"
	"log"
	"sort"

	tezos "github.com/ecadlabs/go-tezos"
)

func (t *TezosListener) retrieveTopNProposal(ctx context.Context, block *tezos.Block, limit int) ([]*tezos.Proposal, error) {
	// Retrieve proposals for the current phase
	proposals, err := t.service.GetProposals(ctx, t.config.GetChainID(), fmt.Sprintf("%d", block.Header.Level-1))

	if err != nil {
		return nil, err
	}
	sort.Slice(proposals, func(i, j int) bool {
		return proposals[i].SupporterCount > proposals[j].SupporterCount
	})

	topN := []*tezos.Proposal{}

	if len(proposals) > 0 {
		for i := range proposals {
			if i+1 > limit {
				return topN, nil
			}

			topN = append(topN, proposals[i])
		}
	}

	return topN, nil
}

func (t *TezosListener) lookForWinningProposal(ctx context.Context, block *tezos.Block) error {
	periodLength := 4096 * 8
	isLastBlockOfVotingPeriod := block.Header.Level%periodLength == 0

	if isLastBlockOfVotingPeriod {
		log.Printf("TezosListener: Inspecting block %s for winning proposal.\n", block.Hash)
		// Retrieve proposals for the current phase
		proposals, err := t.retrieveTopNProposal(ctx, block, 1)

		if err != nil {
			return err
		}

		for i := range proposals {
			// Publish winning proposal
			t.winningProposalChan <- proposals[i]
		}
	}

	return nil
}

func (t *TezosListener) lookForProposalSummary(ctx context.Context, block *tezos.Block) error {

	blockPerDay := 60 * 24

	// The 756th block of the day is about 11am EST
	summaryBlock := 756

	isDailyBlock := block.Header.Level%blockPerDay == summaryBlock

	if isDailyBlock {
		log.Printf("TezosListener: Inspecting block %s for proposal summary.\n", block.Hash)
		// Retrieve proposals for the current phase
		proposals, err := t.retrieveTopNProposal(ctx, block, 3)

		if err != nil {
			return err
		}

		for i := range proposals {
			t.proposalSummaryChan <- proposals[i]
		}
	}

	return nil
}
