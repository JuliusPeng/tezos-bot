package service

import (
	"log"

	"github.com/ecadlabs/tezos-bot/models"
)

// ChainListener interface for required methods of a chain listener
type ChainListener interface {
	Start()
	Stop()
	GetNewVotes() chan *models.Ballot
	GetNewProto() chan string
	GetNewProposal() chan *models.Proposal
	GetProposalUpvotes() chan *models.Proposal
}

// VotePublisher interface for required methods of a vote publisher
type VotePublisher interface {
	Publish(vote *models.Ballot) error
	PublishProtoChange(proto string) error
	PublishProposalInjection(proto *models.Proposal) error
}

// Service main service that listen for new vote on a chain and publish them
type Service struct {
	chainListener ChainListener
	votePublisher VotePublisher
	signals       chan bool
}

// New Create a new service
func New(chainListener ChainListener, votePublisher VotePublisher) *Service {
	return &Service{
		chainListener: chainListener,
		votePublisher: votePublisher,
		signals:       make(chan bool),
	}
}

// Start a the service
func (s *Service) Start() {
	cVote := s.chainListener.GetNewVotes()
	cProto := s.chainListener.GetNewProto()
	cProposal := s.chainListener.GetNewProposal()
	cProposalUpvote := s.chainListener.GetProposalUpvotes()
	go func() {
		for {
			select {
			case vote := <-cVote:
				if err := s.votePublisher.Publish(vote); err != nil {
					log.Printf("%v was not able to be sent due to error: %s", vote, err.Error())
				}
			case proto := <-cProto:
				if err := s.votePublisher.PublishProtoChange(proto); err != nil {
					log.Printf("%v was not able to be sent due to error: %s", proto, err.Error())
				}
			case proposal := <-cProposal:
				if err := s.votePublisher.PublishProposalInjection(proposal); err != nil {
					log.Printf("%v was not able to be sent due to error: %s", *proposal, err.Error())
				}
			case _ = <-cProposalUpvote:
				// TODO(simon) add message for proposal upvote
				break
			case _ = <-s.signals:
				s.chainListener.Stop()
				return
			}
		}
	}()
	s.chainListener.Start()
}

// Stop stop the service
func (s *Service) Stop() {
	s.signals <- true
}
