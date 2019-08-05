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
	GetProposalSummary() chan *models.ProposalSummary
	GetWinningProposal() chan *models.ProposalSummary
}

// VotePublisher interface for required methods of a vote publisher
type VotePublisher interface {
	Publish(vote *models.Ballot) error
	PublishProtoChange(proto string) error
	PublishProposalUpvote(proposal *models.Proposal) error
	PublishProposalInjection(proto *models.Proposal) error
	PublishProposalSummary(proposal *models.ProposalSummary) error
	PublishWinningProposalSummary(proposal *models.ProposalSummary) error
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
	go func() {
		for {
			select {
			case vote := <-s.chainListener.GetNewVotes():
				if err := s.votePublisher.Publish(vote); err != nil {
					log.Printf("%v was not able to be sent due to error: %s", vote, err.Error())
				}
			case proto := <-s.chainListener.GetNewProto():
				if err := s.votePublisher.PublishProtoChange(proto); err != nil {
					log.Printf("%v was not able to be sent due to error: %s", proto, err.Error())
				}
			case proposal := <-s.chainListener.GetNewProposal():
				if err := s.votePublisher.PublishProposalInjection(proposal); err != nil {
					log.Printf("%v was not able to be sent due to error: %s", *proposal, err.Error())
				}
			case proposal := <-s.chainListener.GetProposalUpvotes():
				if err := s.votePublisher.PublishProposalUpvote(proposal); err != nil {
					log.Printf("%v was not able to be sent due to error: %s", *proposal, err.Error())
				}
			case summary := <-s.chainListener.GetProposalSummary():
				if err := s.votePublisher.PublishProposalSummary(summary); err != nil {
					log.Printf("%v was not able to be sent due to error: %s", summary, err.Error())
				}
			case winning := <-s.chainListener.GetWinningProposal():
				if err := s.votePublisher.PublishWinningProposalSummary(winning); err != nil {
					log.Printf("%v was not able to be sent due to error: %s", winning, err.Error())
				}
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
