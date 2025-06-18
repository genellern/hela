package migrator

import (
    "database/sql"
    "errors"
    "strings"
    "time"
)

type Driver string

const (
    MySQL    Driver = "mysql"
    Postgres Driver = "postgres"
)

type Config struct {
    DestinationPath string
    Conn            Connection
}

type Action string

const (
    Create Action = "create"
    Alter  Action = "update"
    Drop   Action = "drop"
)

type Migration struct {
    Version     int
    Name        string
    Table       string
    PackageName string
    Fields      []string
    Action      Action
    Raw         string
}

func (migration Migration) processMigration(connection *Connection) error {
    var fields MigrationFields
    var err error
    if migration.Raw != "" {
        // TODO Handle raw query on db connection
    } else {

        for _, f := range migration.Fields {
            fieldDescription := strings.Split(f, "|")
            fieldDDL := strings.SplitN(fieldDescription[0], ":", 2)

            field := migrationField{
                fieldDDL[0],
                connection.dialect.ProcessDDL(fieldDDL[1]),
            }
            for _, ddl := range fieldDescription[1:] {
                field.ddl += " " + connection.dialect.ProcessDDL(ddl)
            }

            fields = append(fields, field)
        }
    }

    switch migration.Action {
    case Create:
        err = connection.dialect.CreateTable(migration.Table, migration.Raw, fields)
    }

    if err != nil {
        return err
    }
    if done, err := connection.dialect.MarkDone(&migration); done || err != nil {
        if !done {
            err = errors.New("Didn't mark the migration as done> " + migration.Name)
        }
    }

    return err
}

type Migrations []Migration

type MigrationCallback func() Migration
type MigrationCallbacksStack []MigrationCallback

type Connection struct {
    dialect Dialect
    conn    *sql.DB
    Driver
    DSN string
}

type Dialect interface {
    ParseFields() string
    ParseFieldsDDL() []string
    GetLatestMigration() (*MigrationRecord, error)
    ProcessDDL(ddl string) string
    CreateTable(tableName string, rawQuery string, fields MigrationFields) error
    MarkDone(migration *Migration) (bool, error)
}

type MigrationRecord struct {
    Name       string    `json:"name,omitempty" `
    Version    int       `json:"version,omitempty"`
    MigratedOn time.Time `json:"migrated___on"`
}

type migrationField struct {
    name string
    ddl  string
}
type MigrationFields []migrationField

type MysqlDialect struct {
    fields     string
    connection *Connection
}

var _ Dialect = &MysqlDialect{}