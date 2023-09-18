package main

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

type SIFlags uint

const (
	SentinelPort = 26379

	/* SentinelInstance flags */
	SIMaster SIFlags = 1 << iota
	SISlave
	SISentinel
)

type SentinelStat struct {
	cfgPath string
	config  *SentinelConfig
}

func NewSentinelStat(cfgPath string) *SentinelStat {
	s := new(SentinelStat)
	s.cfgPath = cfgPath
	s.config = new(SentinelConfig)
	s.config.Port = SentinelPort

	err := s.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	return s
}

type SentinelAddr struct {
}

type SentinelInstance struct {
	flags SIFlags
	name  string
	addr  SentinelAddr
}

func NewSentinelInstance(name string, flags SIFlags, host string, port int, quorum int, master *SentinelInstance) *SentinelInstance {
	inst := new(SentinelInstance)

	return inst
}

func (s *SentinelStat) Run() {

}

func (s *SentinelStat) ParseConfig() error {
	v := viper.New()
	v.SetConfigFile(s.cfgPath)
	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&s.config)
	if err != nil {
		return err
	}

	return nil
}

func (s *SentinelStat) HandleConfiguration() error {
	/* monitor <name> <host> <port> <quorum> */
	if s.config.Monitor.Quorum <= 0 {
		return errors.New("quorum must be 1 or greater")
	}
	Instance = NewSentinelInstance(
		s.config.Monitor.Name,
		SIMaster,
		s.config.Monitor.IP,
		s.config.Monitor.Port,
		s.config.Monitor.Quorum,
		nil,
	)

	return nil
}
