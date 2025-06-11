package migrator

import (
    "database/sql"
    "errors"
    "fmt"
    "github.com/genellern/hela/cmd/utils"
    _ "github.com/go-sql-driver/mysql"
    "strings"
    "time"
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
    Query(queryStr string, args []interface{}) (*sql.Rows, error)
    Exec(queryStr string, args ...any) (sql.Result, error)
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

type MigrationRecord struct {
    name        string
    version     int
    migrated_on time.Time
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
    var migrations []MigrationOptions
    var err error
    connection = localConnection
    latestMigration, err := getMigrationsDone()
    if err != nil {
        return err
    }

    // Extract migrations
    for _, callback := range *m {
        migrations = append(migrations, callback())
    }

    migrations = filterOutDoneMigrations(migrations, &latestMigration)

    // Run migrations
    for _, migration := range migrations {
        err = processMigration(migration)

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

func filterOutDoneMigrations(migrations []MigrationOptions, latestMigration **MigrationRecord) []MigrationOptions {
    if (*latestMigration) == nil {
        return migrations
    }

    migrations = utils.Filter(migrations, func(migration MigrationOptions) bool {
        return migration.Version > (*latestMigration).version
    })

    return migrations
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

func (c *Connection) Exec(queryStr string, args ...any) (sql.Result, error) {
    fmt.Println("Query> " + queryStr)
    return c.conn.Exec(queryStr, args)
}

func (c *Connection) Query(queryStr string, args []interface{}) (*sql.Rows, error) {
    fmt.Println("Query> " + queryStr)
    return c.conn.Query(queryStr, args...)
}

func getMigrationsDone() (*MigrationRecord, error) {
    var migration *MigrationRecord = nil

    // Check if migrations table already exists
    result, err := connection.Query("show tables", make([]interface{}, 0))
    defer result.Close()

    var tables []string
    for result.Next() {
        var val string
        if err := result.Scan(&val); err != nil {
            return migration, err
        }
        tables = append(tables, val)
    }
    var migrationsExist bool = utils.Contains(tables, "migrations")

    if !migrationsExist {
        return nil, nil
    } else {
        migrationResults, _ := connection.Query("SELECT migration_name, version, migrated_on "+
            "FROM migrations ORDER BY version "+
            "LIMIT 1",
            []interface{}{},
        )
        defer migrationResults.Close()

        for migrationResults.Next() {
            err = migrationResults.Scan(&migration.name, &migration.version, &migration.migrated_on)
            if err != nil {
                return migration, err
            }
        }
    }

    return migration, nil
}