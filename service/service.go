package service

import (
	"fmt"

	"github.com/ecadlabs/tezos-bot/models"
)

// ChainListener interface for required methods of a chain listener
type ChainListener interface {
	Start()
	Stop()
	GetNewVotes() chan *models.Ballot
}

// VotePublisher interface for required methods of a vote publisher
type VotePublisher interface {
	Publish(vote *models.Ballot) error
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
	c := s.chainListener.GetNewVotes()
	go func() {
		for {
			select {
			case vote := <-c:
				if err := s.votePublisher.Publish(vote); err != nil {
					fmt.Printf("%v was not able to be sent due to error: %s", vote, err.Error())
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
