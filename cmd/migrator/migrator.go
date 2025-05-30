package migrator

import (
    "database/sql"
)

type Action string
type Number interface {
    ~int | ~int64 | ~float32 | ~float64 | ~uint
}

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
    Version     Number
    Table       string
    PackageName string
    Fields      []string
    Action      Action
    Raw         string
}

type MigrationCallback func() MigrationOptions
type MigrationsStack []MigrationCallback

type MigrationInterface interface {
    GetFields() []string
    TableName() string
}

//
func (m *MigrationOptions) New(table string, fields []string, action Action, version Number) *MigrationOptions {
    m.Table = table
    m.Version = version
    m.Fields = fields
    m.Action = action

    return m
}

func (m *MigrationsStack) Migrate() error {

    var err error
    for _, callback := range *m {
        options := callback()
        err = processMigration(options)

        if err != nil {
            return err
        }
    }

    return nil
}

func processMigration(options MigrationOptions) error {

    return nil
}

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