package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

type Config struct {
	App    appConfig
	Secret secretCfg
	Server serverCfg
	Tracer tracerCfg
	Syslog syslogCfg
	Redis  redisCfg
	// Datastore
	DBMain mariadbCfg // Galera Cluster
	DBLog  mariadbCfg // Standalone
	Mongo  mongodbCfg
	OneId  oneIdCfg
}

type appConfig struct {
	Name       string
	Version    string
	Mode       string
	PrefixPath string
	StorageDir string
	LogLevel   logrus.Level
	Url        string
}

type secretCfg struct {
	EncryptKey     string
	PrivateKey     *rsa.PrivateKey
	PrivateKeyFile string `env:"SECRET_PRIVATE_KEY"`
}

func (c *appConfig) IsDebug() bool {
	return c.LogLevel == logrus.DebugLevel
}

func LoadConfig(file string, version string) *Config {
	// Read configuration
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalln("load config file error:", err.Error())
	}
	privateKeyPath := viper.GetString("secret.private.key")
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		log.Fatal("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	// Set default configuration
	viper.SetDefault("app.name", "MyOwn")
	viper.SetDefault("app.version", version)
	viper.SetDefault("app.mode", EnvDevelopment)
	viper.SetDefault("app.storage", "./storage")
	viper.SetDefault("app.prefix.path", "/api")
	viper.SetDefault("app.log.level", "info")
	viper.SetDefault("app.url", "")

	// Secret
	viper.SetDefault("secret.encrypt.key", "")

	// Server
	viper.SetDefault("server.listen", ServerListen)
	viper.SetDefault("server.port", ServerPort)
	viper.SetDefault("server.timeout.read", ServerTimeoutRead)
	viper.SetDefault("server.timeout.write", ServerTimeoutWrite)
	viper.SetDefault("server.timeout.idle", ServerTimeoutIdle)
	viper.SetDefault("server.header", viper.GetString("app.name"))
	viper.SetDefault("server.proxy.header", fiber.HeaderXForwardedFor)
	viper.SetDefault("server.enable.cors", "false")
	viper.SetDefault("server.buffer.read", ReadBufferSize)
	viper.SetDefault("server.body.limit", 10*1024*1024*1024) // 10GB
	// Syslog
	viper.SetDefault("syslog.enable", "false")
	viper.SetDefault("syslog.server", "127.0.0.1")
	viper.SetDefault("syslog.port", "514")
	viper.SetDefault("syslog.protocol", "udp")

	// DBMain
	viper.SetDefault("db.main.host", MariadbHost)
	viper.SetDefault("db.main.port", MariadbPort)
	viper.SetDefault("db.main.username", "")
	viper.SetDefault("db.main.password", "")
	viper.SetDefault("db.main.database", "")
	viper.SetDefault("db.main.migration", false)
	// DBLog
	viper.SetDefault("db.log.host", MariadbHost)
	viper.SetDefault("db.log.port", MariadbPort)
	viper.SetDefault("db.log.username", "")
	viper.SetDefault("db.log.password", "")
	viper.SetDefault("db.log.database", "")
	// Redis
	viper.SetDefault("redis.host", RedisHost)
	viper.SetDefault("redis.port", RedisPort)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.database", RedisDb)

	// Oneid
	viper.SetDefault("oneid.client_id", "")
	viper.SetDefault("oneid.client_secret", "")
	viper.SetDefault("oneid.ref_code", "")
	viper.SetDefault("oneid.redirect_url", "")
	viper.SetDefault("oneid.url", "")

	config := &Config{
		App: appConfig{
			Name:       viper.GetString("app.name"),
			Version:    viper.GetString("app.version"),
			Mode:       viper.GetString("app.mode"),
			PrefixPath: viper.GetString("app.prefix.path"),
			StorageDir: viper.GetString("app.storage"),
			Url:        viper.GetString("app.url"),
		},
		Secret: secretCfg{
			EncryptKey:     viper.GetString("secret.encrypt.key"),
			PrivateKeyFile: privateKeyPath,
			PrivateKey:     privateKey,
		},
		Server: serverCfg{
			ListenIp:       viper.GetString("server.listen"),
			Port:           viper.GetString("server.port"),
			TimeoutRead:    util.ParseDuration(viper.GetString("server.timeout.read")),
			TimeoutWrite:   util.ParseDuration(viper.GetString("server.timeout.write")),
			TimeoutIdle:    util.ParseDuration(viper.GetString("server.timeout.idle")),
			ServerHeader:   viper.GetString("server.header"),
			ProxyHeader:    viper.GetString("server.proxy.header"),
			EnableCORS:     viper.GetBool("server.enable.cors"),
			ReadBufferSize: viper.GetInt("server.buffer.read"),
			BodyLimit:      viper.GetInt("server.body_limit"),
		},
		Syslog: syslogCfg{
			Enable:   viper.GetString("syslog.enable") == "true",
			Server:   viper.GetString("syslog.server"),
			Port:     viper.GetString("syslog.port"),
			Protocol: viper.GetString("syslog.protocol"),
		},
		Tracer: tracerCfg{
			Enable: viper.GetString("tracer.enable") == "true",
			Url:    viper.GetString("tracer.url"),
		},
		DBMain: mariadbCfg{
			Host:      viper.GetString("db.main.host"),
			Port:      viper.GetString("db.main.port"),
			User:      viper.GetString("db.main.username"),
			Password:  viper.GetString("db.main.password"),
			Database:  viper.GetString("db.main.database"),
			Migration: viper.GetBool("db.main.migration"),
		},
		DBLog: mariadbCfg{
			Host:     viper.GetString("db.log.host"),
			Port:     viper.GetString("db.log.port"),
			User:     viper.GetString("db.log.username"),
			Password: viper.GetString("db.log.password"),
			Database: viper.GetString("db.log.database"),
		},
		Mongo: mongodbCfg{
			Connection: viper.GetString("mongo.connection"),
			Database:   viper.GetString("mongo.database"),
		},
		Redis: redisCfg{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetString("redis.port"),
			Password: viper.GetString("redis.password"),
			Database: viper.GetString("redis.database"),
		},
		OneId: oneIdCfg{
			ClientId:     viper.GetString("oneid.client_id"),
			ClientSecret: viper.GetString("oneid.client_secret"),
			RefCode:      viper.GetString("oneid.ref_code"),
			RedirectUrl:  viper.GetString("oneid.redirect_url"),
			Url:          viper.GetString("oneid.url"),
			IalLevel:     viper.GetFloat64("oneid.ial_level"),
			ForgotUrl:    viper.GetString("oneid.forgot_url"),
		},
	}
	config.App.LogLevel, err = logrus.ParseLevel(viper.GetString("app.log.level"))
	if err != nil {
		config.App.LogLevel = logrus.InfoLevel
	}
	return config
	// return &Config{}
}
