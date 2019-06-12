package main

import (
	"flag"
	"log"

	"github.com/ecadlabs/tezos-bot/config"
	"github.com/ecadlabs/tezos-bot/listen"
	"github.com/ecadlabs/tezos-bot/publish"
	"github.com/ecadlabs/tezos-bot/service"
)

func main() {

	c := config.Config{
		RPCURL:               "https://mainnet-node.tzscan.io",
		ChainID:              "main",
		RetryCount:           100,
		History:              false,
		MonitorVote:          true,
		MonitorProtocol:      true,
		HistoryStartingBlock: 0,
	}

	var configFile string
	flag.StringVar(&configFile, "c", "./config.yaml", "Config file.")
	flag.Parse()
	c.Load(configFile)

	l, err := listen.NewTezosListener(c)

	if err != nil {
		log.Printf(err.Error())
		return
	}

	var p service.VotePublisher

	if c.GetTwitterAccessToken() == "" {
		log.Printf("Twitter access token not configured posting vote to stdout\n")
		p = &publish.DebugPublisher{}
	} else {
		p, err = publish.NewTwitterPublisher(c)

		if err != nil {
			log.Printf(err.Error())
			return
		}
	}

	s := service.New(l, p)

	log.Println("Bot started...")

	s.Start()
}
