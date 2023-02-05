package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

// PoolConfig contains the connection pool related properties.
type PoolConfig struct {
	MaxOpenConnections        int
	MaxIdleConnections        int
	ConnectionMaxLifetimeSecs int
	ConnectionMaxIdleTimeSecs int
}

// Config represents the postgres config.
type Config struct {
	Pool               *PoolConfig
	Host               string
	Port               string
	Password           string
	Database           string
	SSL                string
	AppName            string
	Username           string
	ConnectTimeoutSecs int
}

// connString represents a connection string type.
type connString struct {
	dsn string
}

// asDSN constructs a DSN connection string out of a regular config.
func (c *Config) asDSN() (*connString, error) {
	if c.Host == "" {
		return nil, errors.New("host can`t be empty")
	}

	if c.Port == "" {
		return nil, errors.New("port can`t be empty")
	}

	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 255 || port > 65535 {
		return nil, errors.New("port is invalid")
	}

	if c.Username == "" {
		return nil, errors.New("username can`t be empty")
	}

	if c.SSL == "" {
		return nil, errors.New("ssl mode can`t be empty")
	}

	if c.ConnectTimeoutSecs == 0 {
		return nil, errors.New("connect timeout seconds can`t be empty")
	}

	dsn := fmt.Sprintf("host=%s port=%v user=%s dbname=%s sslmode=%s fallback_application_name=%s connect_timeout=%d",
		c.Host,
		port,
		c.Username,
		c.Database,
		c.SSL,
		c.AppName,
		c.ConnectTimeoutSecs,
	)

	if c.Password != "" {
		dsn = fmt.Sprintf("%s password=%s", dsn, c.Password)
	}

	return &connString{dsn: dsn}, nil
}

// clientOptionsInternal represents an internal options type for the client.
type clientOptionsInternal struct {
	pool *PoolConfig
	dsn  string
}

// ClientOptions represents the client`s options.
type ClientOptions struct {
	internal     *clientOptionsInternal
	databaseName string
}

// PostgresClient represents a structure containing all postgres client related properties.
type PostgresClient struct {
	db   *sql.DB
	opts *ClientOptions
}

// setPool applies the PoolConfig settings to a PostgresClient.
func (pc *PostgresClient) setPool() {
	if pc.opts.internal.pool == nil {
		return
	}

	pool := pc.opts.internal.pool
	pc.db.SetMaxOpenConns(pool.MaxOpenConnections)
	pc.db.SetMaxIdleConns(pool.MaxIdleConnections)
	pc.db.SetConnMaxLifetime(time.Duration(pool.ConnectionMaxLifetimeSecs) * time.Second)
	pc.db.SetConnMaxIdleTime(time.Duration(pool.ConnectionMaxIdleTimeSecs) * time.Second)
}

// Database returns the internal PostgresClient database.
func (pc *PostgresClient) Database() *sql.DB {
	return pc.db
}

// Connect initializes a connection to the database.
func (pc *PostgresClient) Connect(ctx context.Context) error {
	c, err := sql.Open("postgres", pc.opts.internal.dsn)
	if err != nil {
		return err
	}

	pc.db = c
	pc.setPool()

	return nil
}

// Ping performs a ping to the database host.
func (pc *PostgresClient) Ping(ctx context.Context) error {
	return pc.db.Ping()
}

// NewPostgresClient creates a new PostgresClient.
func NewPostgresClient(ctx context.Context, c *Config) (*PostgresClient, error) {
	cs, err := c.asDSN()
	if err != nil {
		return nil, err
	}

	clientInternalOpts := clientOptionsInternal{dsn: cs.dsn}

	if c.Pool != nil {
		clientInternalOpts.pool = c.Pool
	}

	return &PostgresClient{
		opts: &ClientOptions{
			internal:     &clientInternalOpts,
			databaseName: c.Database,
		},
	}, nil
}
