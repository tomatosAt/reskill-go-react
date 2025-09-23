package config

type syslogCfg struct {
	Enable   bool
	Server   string
	Port     string
	Protocol string
}
