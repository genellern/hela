package migrator

import (
    "database/sql"
)
type Action string

const (
    Create Action = "create"
    Alter  Action = "update"
    Drop   Action = "drop"
)

type Dialect string

const (
    MySQL Dialect = "MySQL"
)

type Connection struct {
    Dialect Dialect
    DSN     string
    conn    *sql.DB
}

type ConnectionInterface interface {
    Open() error
    Close() error
}

type Config struct {
    DestinationPath string
    Conn            ConnectionInterface
}

type MigrationOptions struct {
    Table       string
    PackageName string
    Fields      []string
    Action      Action
}

type MigrationInterface interface {
    GetFields() []string
    TableName() string
}

func (m MigrationOptions) GetFields() []string {
    return m.Fields

// Connection

func (c *Connection) Open() error {
    db, err := sql.Open(string(c.Dialect), c.DSN)
    if err != nil {
        panic(err)
    }
    c.conn = db
    return nil
}

func (c *Connection) Close() error {
    return c.conn.Close()
}