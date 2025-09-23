package config

import "time"

type serverCfg struct {
	ListenIp       string
	Port           string
	TimeoutRead    time.Duration
	TimeoutWrite   time.Duration
	TimeoutIdle    time.Duration
	ServerHeader   string
	ProxyHeader    string
	EnableCORS     bool
	ReadBufferSize int
	BodyLimit      int
}
