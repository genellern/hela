package migrator

import (
    "database/sql"
    "fmt"
)

func (c *Connection) Open() error {
    db, err := sql.Open(string(c.Driver), c.DSN)

    if err != nil {
        return err
    }
    c.conn = db
    c.SetDialect()
    return nil
}

func (c *Connection) Close() error {
    return c.conn.Close()
}

func (c *Connection) Exec(queryStr string, args ...any) (sql.Result, error) {
    fmt.Println("Exec> " + queryStr)
    return c.conn.Exec(queryStr, args)
}

func (c *Connection) Query(queryStr string, args ...any) (*sql.Rows, error) {
    fmt.Println("Query> " + queryStr)
    return c.conn.Query(queryStr, args)
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