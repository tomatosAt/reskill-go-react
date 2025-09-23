package database

import (
	"database/sql"
	"fmt"

	// "github.com/carlescere/scheduler"
	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Host      string
	Port      string
	Username  string
	Password  string
	Name      string
	DebugMode bool
	Migration bool
}

type Client struct {
	config     Config
	dsn        string
	ctx        *gorm.DB
	sql        *sql.DB
	gormConfig gorm.Config
	// job        *scheduler.Job
	logger *logrus.Entry
}

func NewWithConfig(cfg Config, logger *logrus.Logger) Client {
	return Client{
		config: cfg,
		logger: logger.WithField("package", "database"),
	}
}

func (c *Client) Connect() error {
	return c.ConnectWithGormConfig(gorm.Config{})
}

func (c *Client) ConnectWithGormConfig(gormCfg gorm.Config) error {
	_ = c.Close()
	c.gormConfig = gormCfg
	c.dsn = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		c.config.Username,
		c.config.Password,
		c.config.Host,
		c.config.Port,
		c.config.Name,
		"utf8mb4",
	)
	var err error
	c.ctx, err = gorm.Open(mysql.Open(c.dsn), &gormCfg)
	if err != nil {
		return err
	}

	c.sql, err = c.ctx.DB()
	if err != nil {
		return err
	}

	if c.config.DebugMode {
		c.ctx = c.ctx.Debug()
	}
	if c.config.Migration {
		if err := c.ctx.AutoMigrate(
			model.Game{},
			model.User{},
			model.Cart{},
			model.Purchase{},
			model.Transaction{},
			model.TransactionAPILog{},
			model.PreRegister{},
		); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Ctx() *gorm.DB {
	return c.ctx
}

func (c *Client) SqlDB() *sql.DB {
	return c.sql
}

func (c *Client) MigrateDatabase(tables []interface{}) error {
	tx := c.ctx.Begin()
	for _, t := range tables {
		if err := tx.AutoMigrate(t); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (c *Client) Close() error {
	// do nothing if no connection
	if c.sql == nil {
		return nil
	}
	if err := c.sql.Close(); err != nil {
		return err
	}
	return nil
}
