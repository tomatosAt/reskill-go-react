package config

type mariadbCfg struct {
	Host      string
	Port      string
	User      string
	Password  string
	Database  string
	Migration bool
}

type mongodbCfg struct {
	Connection string
	Database   string
}

type redisCfg struct {
	Host     string
	Port     string
	Password string
	Database string
}
