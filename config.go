package main

type MonitorConfig struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Quorum int    `yaml:"quorum"`
}

type SentinelConfig struct {
	Port    int           `yaml:"port"`
	Monitor MonitorConfig `yaml:"monitor"`
}
