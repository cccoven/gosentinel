package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	"net"
	"time"
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
	client *redis.Client
	flags  SIFlags
	name   string
	addr   *SentinelAddr

	quorum    int
	Sentinels map[string]*SentinelInstance
	Slaves    map[string]*SentinelInstance
}

func NewSentinelInstance(name string, flags SIFlags, host string, port int, quorum int, master *SentinelInstance) *SentinelInstance {
	inst := new(SentinelInstance)
	inst.flags = flags
	inst.name = name
	inst.addr = NewSentinelAddr(host, port)
	inst.quorum = quorum
	inst.Sentinels = make(map[string]*SentinelInstance)
	inst.Slaves = make(map[string]*SentinelInstance)

	inst.client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", inst.addr.IP, inst.addr.Port),
	})

	if flags&SIMaster != 0 {
		Sentinel.Masters[inst.name] = inst
	} else if flags&SISlave != 0 {
		inst.Slaves[inst.name] = inst
	} else if flags&SISentinel != 0 {
		inst.Sentinels[inst.name] = inst
	}

	return inst
}

func (s *SentinelStat) Run() {
	t := time.NewTicker(time.Second)
	for _ = range t.C {
		s.run()
	}
}

func (s *SentinelStat) run() {
	s.HandleInstances(s.Masters)
}

func (s *SentinelStat) HandleInstances(instances map[string]*SentinelInstance) {
	for _, inst := range instances {
		inst.Handle()
		if inst.flags&SIMaster != 0 {
			s.HandleInstances(inst.Slaves)
			s.HandleInstances(inst.Sentinels)
		}
	}
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

func (si *SentinelInstance) Handle() {
	si.SendPeriodCommands()
}

func (si *SentinelInstance) SendPeriodCommands() {
	if si.flags&SISentinel == 0 {
		infoResult, _ := si.client.Info(context.Background()).Result()
		fmt.Println(infoResult)
	}
}

func (si *SentinelInstance) CheckSubjectivelyDown() {

}

func (si *SentinelInstance) CheckObjectivelyDown() {

}
