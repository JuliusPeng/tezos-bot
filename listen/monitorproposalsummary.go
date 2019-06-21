package listen

import (
	"context"
	"fmt"
	"log"
	"sort"

	tezos "github.com/ecadlabs/go-tezos"
	"github.com/ecadlabs/tezos-bot/models"
)

func (t *TezosListener) retrieveTopNProposal(ctx context.Context, level int, limit int) ([]*tezos.Proposal, error) {
	// Retrieve proposals for the current phase
	proposals, err := t.service.GetProposals(ctx, t.config.GetChainID(), fmt.Sprintf("%d", level-1))

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
		proposals, err := t.retrieveTopNProposal(ctx, block.Header.Level, 1)

		if err != nil {
			return err
		}

		for i := range proposals {
			// Publish winning proposal
			t.winningProposalChan <- &models.ProposalSummary{
				Proposal: *proposals[i],
				Cycle:    (block.Header.Level / 4096) - 1,
			}
		}
	}

	return nil
}

func (t *TezosListener) lookForProposalSummary(ctx context.Context, block *tezos.Block) error {
	if block.Header.Level%4096 == 0 {
		log.Printf("TezosListener: Inspecting block %s for proposal summary.\n", block.Hash)
		// Retrieve proposals for the current phase
		proposals, err := t.retrieveTopNProposal(ctx, block.Header.Level, 3)

		if err != nil {
			return err
		}

		previousProposals, err := t.service.GetProposals(ctx, t.config.GetChainID(), fmt.Sprintf("%d", block.Header.Level-4096))

		if err != nil {
			return err
		}

		for i := range proposals {
			previousSupporters := 0
			for _, prev := range previousProposals {
				if prev.ProposalHash == proposals[i].ProposalHash {
					previousSupporters = prev.SupporterCount
				}
			}

			newSupporters := proposals[i].SupporterCount - previousSupporters

			t.proposalSummaryChan <- &models.ProposalSummary{
				Proposal:      *proposals[i],
				Cycle:         (block.Header.Level / 4096) - 1,
				NewSupporters: newSupporters,
			}
		}
	}

	return nil
}
