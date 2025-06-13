package migrator

import (
    "database/sql"
    "fmt"
)

func (c *Connection) Open() error {
    db, err := sql.Open(string(c.dialectStr), c.DSN)
    if err != nil {
        panic(err)
    }
    c.conn = db
    return nil
}

func (c *Connection) Close() error {
    return c.conn.Close()
}

func (c *Connection) Exec(queryStr string, args ...any) (sql.Result, error) {
    fmt.Println("Query> " + queryStr)
    return c.conn.Exec(queryStr, args)
}

func (c *Connection) Query(queryStr string, args []interface{}) (*sql.Rows, error) {
    return c.dialectConcrete.Query(queryStr, args)
}