package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config struct containing all configurable parameter for the tezos bot
type Config struct {
	RPCURL                   string `yaml:"rpc_url"`
	TwitterAccessToken       string `yaml:"twitter_access_token"`
	TwitterAccessTokenSecret string `yaml:"twitter_access_token_secret"`
	TwitterConsummerID       string `yaml:"twitter_consummer_id"`
	TwitterConsummerKey      string `yaml:"twitter_consummer_key"`
	ChainID                  string `yaml:"chain"`
	RetryCount               int    `yaml:"retry_count"`
	History                  bool   `yaml:"history"`
}

// GetRPCURL returns the tezos rpc endpoint
func (c Config) GetRPCURL() string {
	return c.RPCURL
}

// IsHistory returns true if listener should read history
func (c Config) IsHistory() bool {
	return c.History
}

// GetRetryCount returns the max retry attempt to connect to tezos node
func (c Config) GetRetryCount() int {
	return c.RetryCount
}

// GetChainID returns the chain ID
func (c Config) GetChainID() string {
	return c.ChainID
}

// GetTwitterAccessToken returns the twitter access token
func (c Config) GetTwitterAccessToken() string {
	return c.TwitterAccessToken
}

// GetTwitterAccessTokenSecret returns the twitter access token secret
func (c Config) GetTwitterAccessTokenSecret() string {
	return c.TwitterAccessTokenSecret
}

// GetTwitterConsummerID returns the twitter consummer id
func (c Config) GetTwitterConsummerID() string {
	return c.TwitterConsummerID
}

// GetTwitterConsummerKey returns the twitter consummer key
func (c Config) GetTwitterConsummerKey() string {
	return c.TwitterConsummerKey
}

// Load read a config file and unmarshal it into the config struct
func (c *Config) Load(name string) error {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(buf, c); err != nil {
		return err
	}

	return nil
}
