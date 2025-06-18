package migrator

import (
    "database/sql"
    "errors"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

func (c *Connection) Open() error {
    fmt.Println("Opening connection")

    if c == nil {
        return errors.New("connection is nil")
    }
    if c.Driver == "" {
        return errors.New("Driver is empty")
    }
    if c.DSN == "" {
        return errors.New("DSN is empty")
    }

    db, err := sql.Open(string(c.Driver), c.DSN)

    if err != nil {
        return err
    }

    c.conn = db
    c.SetDialect()
    return nil
}

func (c *Connection) Close() error {
    fmt.Println("Closing connection")

    return c.conn.Close()
}

func (c *Connection) Exec(queryStr string) (sql.Result, error) {
    fmt.Println("Exec> " + queryStr)
    return c.conn.Exec(queryStr)
}

func (c *Connection) Query(queryStr string) (*sql.Rows, error) {
    fmt.Println("Query> " + queryStr)
    return c.conn.Query(queryStr)
}

func (c *Connection) SetDialect() {
    switch c.Driver {
    case MySQL:
        builder := &MySQLDialectBuilder{}
        c.dialect = builder.Build(c)
    }
}

func (c *Connection) GetLatestMigration() (*MigrationRecord, error) {
    return c.dialect.GetLatestMigration()
}