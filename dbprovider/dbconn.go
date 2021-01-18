package dbprovider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TxKey struct{}

var ErrNilTx = errors.New("tx is nil")

type DBProvider interface {
	ExecuteTx(ctx context.Context, fn func(txCtx context.Context) error) error
	Begin(ctx context.Context) (txCtx context.Context)
	Commit(txCtx context.Context) error
	Rollback(txCtx context.Context) error
	DB() (*sql.DB, error)
}

// DBType enumerate all database type
type DBType string

func (dbType DBType) String() string {
	return string(dbType)
}

const (
	// Mysql mysql db
	Mysql DBType = "mysql"
	// Pg postgres db
	Pg DBType = "postgres"
	// Sqlite sqlite db
	Sqlite DBType = "sqlite"
)

type DBConfig struct {
	Host       string `yaml:"host"`
	Port       int32  `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	DBName     string `yaml:"dbname"`
	SearchPath string `yaml:"search_path"` // pg should setting this value. It will restrict access to db schema.
	InMemory   bool   `yaml:"in_memory"`   // sqlite setting
	SSLMode    bool   `yaml:"ssl_mode"`
	Type       DBType `yaml:"type"`
	PoolConfig
}

// Addr return db connection addresss
func (config DBConfig) Addr() string {
	return config.Host + ":" + strconv.Itoa(int(config.Port))
}

// pgSSLMode return pg ssl mode string setting
func (config DBConfig) pgSSLMode() string {
	if config.SSLMode {
		return "require"
	}
	return "disable"
}

// ConnString return connection string based db type
func (config DBConfig) ConnString() string {
	var connectionString string

	switch config.Type {
	case Mysql:
		connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true",
			config.User, config.Password, config.Addr(), config.DBName)
	case Pg:
		connectionString = fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s`,
			config.User, config.Password, config.Host, config.Port, config.DBName, config.pgSSLMode())
		if strings.TrimSpace(config.SearchPath) != "" {
			connectionString = fmt.Sprintf("%s search_path=%s", connectionString, config.SearchPath)
		}
	case Sqlite:
		connectionString = fmt.Sprintf("%s", config.Host)
		if config.InMemory {
			connectionString = fmt.Sprintf("%s:memory:?cache=shared", connectionString)
		}
	}
	return connectionString
}

type PoolConfig struct {
	MaxIdleConns   int `yaml:"maxidleconns"`
	MaxOpenConns   int `yaml:"maxopenconns"`
	MaxIdletimeSec int `yaml:"maxidletimesec"`
	MaxLifetimeSec int `yaml:"maxlifetimesec"`
}

func (c PoolConfig) ConfigWithDB(db *sql.DB) {
	if c.MaxIdleConns != 0 {
		db.SetMaxIdleConns(c.MaxIdleConns)
	}
	if c.MaxOpenConns != 0 {
		db.SetMaxOpenConns(c.MaxOpenConns)
	}
	if c.MaxLifetimeSec != 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(c.MaxLifetimeSec))
	}
	if c.MaxIdletimeSec != 0 {
		db.SetConnMaxIdleTime(time.Second * time.Duration(c.MaxIdletimeSec))
	}
}
