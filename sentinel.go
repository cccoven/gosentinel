package main

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"net"
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
	ConfigPath string
	Config     *SentinelConfig
	Masters    map[string]*SentinelInstance
}

func NewSentinelStat(cfgPath string) *SentinelStat {
	s := new(SentinelStat)
	s.ConfigPath = cfgPath
	s.Config = new(SentinelConfig)
	s.Config.Port = SentinelPort
	s.Masters = make(map[string]*SentinelInstance)

	err := s.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	return s
}

type SentinelAddr struct {
	Host string
	IP   string
	Port int
}

func NewSentinelAddr(host string, port int) *SentinelAddr {
	sa := new(SentinelAddr)

	addrs, _ := net.LookupHost(host)
	sa.Host = host
	sa.IP = addrs[0]
	sa.Port = port

	return sa
}

type SentinelInstance struct {
	flags SIFlags
	name  string
	addr  *SentinelAddr

	quorum int
}

func NewSentinelInstance(name string, flags SIFlags, host string, port int, quorum int, master *SentinelInstance) *SentinelInstance {
	inst := new(SentinelInstance)
	inst.flags = flags
	inst.name = name
	inst.addr = NewSentinelAddr(host, port)
	inst.quorum = quorum

	return inst
}

func (s *SentinelStat) Run() {

}

func (s *SentinelStat) ParseConfig() error {
	v := viper.New()
	v.SetConfigFile(s.ConfigPath)
	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&s.Config)
	if err != nil {
		return err
	}

	return nil
}

func (s *SentinelStat) HandleConfiguration() error {
	/* monitor <name> <host> <port> <quorum> */
	if s.Config.Monitor.Quorum <= 0 {
		return errors.New("quorum must be 1 or greater")
	}
	Instance = NewSentinelInstance(
		s.Config.Monitor.Name,
		SIMaster,
		s.Config.Monitor.Host,
		s.Config.Monitor.Port,
		s.Config.Monitor.Quorum,
		nil,
	)

	return nil
}
