package migrator

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "strings"
)

type Action string

const (
    Create Action = "create"
    Alter  Action = "update"
    Drop   Action = "drop"
)

type Dialect string

const (
    MySQL Dialect = "mysql"
)

type Connection struct {
    Dialect Dialect
    // username:password@protocol(address)/dbname?param=value
    DSN  string
    conn *sql.DB
}

type ConnectionInterface interface {
    Open() error
    Close() error
    Query(queryStr string, args []interface{}) (sql.Result, error)
}

type Config struct {
    DestinationPath string
    Conn            ConnectionInterface
}

var connection ConnectionInterface

type MigrationOptions struct {
    Version     int
    Name        string
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
func (m *MigrationOptions) New(table string, fields []string, action Action, version int) *MigrationOptions {
    m.Table = table
    m.Version = version
    m.Fields = fields
    m.Action = action

    return m
}

func (m *MigrationsStack) Migrate(localConnection ConnectionInterface) error {
    connection = localConnection
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

type field struct {
    name string
    ddl  string
}
type fields []field

func (f fields) ddl() []string {
    ddl := []string{}
    for _, field := range f {
        ddl = append(ddl, field.name+" "+field.ddl)
    }

    return ddl
}

func processMigration(options MigrationOptions) error {
    var fields fields
    var err error
    if options.Raw != "" {
        // TODO Handle raw query on db connection
    } else {

        for _, f := range options.Fields {
            fieldDescription := strings.Split(f, "|")
            fieldDDL := strings.SplitN(fieldDescription[0], ":", 2)

            field := field{
                fieldDDL[0],
                processDDL(fieldDDL[1]),
            }
            for _, ddl := range fieldDescription[1:] {
                field.ddl += " " + processDDL(ddl)
            }

            fields = append(fields, field)
        }
    }

    switch options.Action {
    case Create:
        err = createTable(options.Table, options.Raw, fields)
    }

    return err
}

func createTable(table, rawQuery string, fields fields) error {

    var err error
    if rawQuery == "" {
        rawQuery = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ("+
            "%s"+
            ")", table, strings.Join(fields.ddl(), ",\n"))
    }

    _, err = connection.Query(rawQuery, []interface{}{})
    fmt.Println(rawQuery)

    return err
}

func processDDL(ddl string) string {
    switch ddl {
    case "nullable":
        return "DEFAULT NULL"
    case "uuid":
        return "CHAR(36)"
    case "required":
        return "NOT NULL"
    case "unique":
        // TODO add indexes to DDL
        return ""
    case "boolean":
        return "tinyint(2)"
    default:
        return ddl
    }
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

func (c *Connection) Query(queryStr string, args []interface{}) (sql.Result, error) {
    fmt.Println("Query> " + queryStr)
    return c.conn.Exec(queryStr, args...)
}