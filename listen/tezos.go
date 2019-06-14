package listen

import (
	"context"
	"log"
	"net/http"

	"github.com/ecadlabs/go-tezos"
	"github.com/ecadlabs/tezos-bot/models"
)

var lastBlock *tezos.Block

// TezosConfig interface with method necessary to obtain tezos listener configurable parameter
type TezosConfig interface {
	GetRPCURL() string
	GetChainID() string
	GetRetryCount() int
	IsMonitorVote() bool
	IsMonitorProtocol() bool
	IsMonitorProposal() bool
	IsHistory() bool
	GetHistoryStartingBlock() int
}

// TezosListener is a struct containing information necessary to monitor the tezos chain
type TezosListener struct {
	service             *tezos.Service
	votesChan           chan *models.Ballot
	protoChan           chan string
	newProposalChan     chan *models.Proposal
	proposalUpvoteChan  chan *models.Proposal
	proposalSummaryChan chan *models.ProposalSummary
	winningProposalChan chan *models.ProposalSummary
	cache               *cache
	signals             chan bool
	config              TezosConfig
	bStreaming          BlockStreamingFunc
}

// NewTezosListener create a new TezosListener
func NewTezosListener(config TezosConfig) (*TezosListener, error) {
	client, err := tezos.NewRPCClient(http.DefaultClient, config.GetRPCURL())
	if err != nil {
		return nil, err
	}

	var bStreamingFunc BlockStreamingFunc = MonitorBlockStreamingFunc

	if config.IsHistory() {
		bStreamingFunc = HistoryBlockStreamingFunc
	}

	return &TezosListener{
		service:             &tezos.Service{Client: client},
		cache:               newCache(),
		votesChan:           make(chan *models.Ballot),
		protoChan:           make(chan string),
		newProposalChan:     make(chan *models.Proposal),
		proposalUpvoteChan:  make(chan *models.Proposal),
		proposalSummaryChan: make(chan *models.ProposalSummary),
		winningProposalChan: make(chan *models.ProposalSummary),
		signals:             make(chan bool),
		config:              config,
		bStreaming:          bStreamingFunc,
	}, nil
}

// Start start monitoring the chain and push new Ballot in the ballot channel
func (t *TezosListener) Start() {
	ctx := context.Background()
	cBlockHash := make(chan string)
	defer close(cBlockHash)
	go func() {
		for {
			select {
			case _ = <-t.signals:
				return
			default:
				err := t.bStreaming(ctx, t.config, t.service, cBlockHash)
				if err != nil {
					panic(err.Error())
				}
				return
			}
		}
	}()

	for hash := range cBlockHash {
		// cBlockHash channel can emit the same has multiple time
		// In order to avoid duplicate we check if it has already been processed
		if !t.cache.Has(hash) {
			t.cache.Add(hash)
			block, err := t.service.GetBlock(ctx, t.config.GetChainID(), hash)

			if err != nil {
				log.Printf("Block: %s skipped because of error: %s\n", hash, err.Error())
				continue
			}

			periodKind, err := t.service.GetCurrentPeriodKind(ctx, t.config.GetChainID(), block.Hash)

			if err != nil {
				log.Printf("Block: %s skipped because of error: %s\n", hash, err.Error())
				continue
			}

			if t.config.IsMonitorVote() && (periodKind.IsTestingVote() || periodKind.IsPromotionVote()) {
				err = t.lookForBallot(ctx, block, periodKind)
				if err != nil {
					log.Printf("Block: %s skipped because of error: %s\n", hash, err.Error())
					continue
				}
			}

			if t.config.IsMonitorProtocol() {
				err = t.lookForProtocolChange(ctx, block)
				if err != nil {
					log.Printf("Block: %s skipped because of error: %s\n", hash, err.Error())
					continue
				}
			}

			if t.config.IsMonitorProposal() && periodKind.IsProposal() {
				err = t.lookForProposal(ctx, block)
				if err != nil {
					log.Printf("Block: %s skipped because of error: %s\n", hash, err.Error())
					continue
				}
			}

			if t.config.IsMonitorProposal() && (periodKind.IsProposal()) {
				err = t.lookForProposalSummary(ctx, block)
				if err != nil {
					log.Printf("Block: %s skipped because of error: %s\n", hash, err.Error())
					continue
				}
			}

			if t.config.IsMonitorProposal() && (periodKind.IsTestingVote()) {
				err = t.lookForWinningProposal(ctx, block)
				if err != nil {
					log.Printf("Block: %s skipped because of error: %s\n", hash, err.Error())
					continue
				}
			}
		}
	}
}

// Stop stop the tezos listener
func (t *TezosListener) Stop() {
	t.signals <- true
}

// GetNewVotes returns a Ballot channel
func (t *TezosListener) GetNewVotes() chan *models.Ballot {
	return t.votesChan
}

// GetNewProto returns a proto channel
func (t *TezosListener) GetNewProto() chan string {
	return t.protoChan
}

// GetNewProposal returns a proto channel
func (t *TezosListener) GetNewProposal() chan *models.Proposal {
	return t.newProposalChan
}

// GetProposalUpvotes returns a proto channel
func (t *TezosListener) GetProposalUpvotes() chan *models.Proposal {
	return t.proposalUpvoteChan
}

// GetProposalSummary returns a proposal summary channel
func (t *TezosListener) GetProposalSummary() chan *models.ProposalSummary {
	return t.proposalSummaryChan
}

// GetWinningProposal returns a proposal summary channel
func (t *TezosListener) GetWinningProposal() chan *models.ProposalSummary {
	return t.winningProposalChan
}
