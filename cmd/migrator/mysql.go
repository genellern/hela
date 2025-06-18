package migrator

import (
    "fmt"
    "github.com/genellern/hela/cmd/utils"
    "strings"
    "time"
)

func (d *MysqlDialect) ParseFields() string {
    return d.fields
}

func (d *MysqlDialect) ParseFieldsDDL() []string {
    return strings.Split(d.fields, " ")
}

type DialectBuilderInterface interface {
    Build(c *Connection) Dialect
}

type MySQLDialectBuilder struct {
}

func (d *MySQLDialectBuilder) Build(c *Connection) Dialect {
    return &MysqlDialect{
        connection: c,
    }
}

func (d *MysqlDialect) ProcessDDL(ddl string) string {
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

func (d *MysqlDialect) GetLatestMigration() (*MigrationRecord, error) {
    var migration = &MigrationRecord{}

    // Check if migrations table already exists
    result, err := d.connection.Query("show tables")
    defer result.Close()
    if err != nil {
        return nil, err
    }

    var tables []string
    for result.Next() {
        var val string
        if err := result.Scan(&val); err != nil {
            return nil, err
        }
        tables = append(tables, val)
    }

    var migrationsExist bool = utils.Contains(tables, "migrations")
    if !migrationsExist {
        return nil, nil
    } else {
        migrationResults, err := d.connection.Query("SELECT migration_name, migration_version, migrated_on " +
            "FROM migrations ORDER BY migration_version DESC " +
            "LIMIT 1",
        )
        if err != nil {
            return nil, err
        }

        defer migrationResults.Close()

        for migrationResults.Next() {
            err = migrationResults.Scan(&migration.Name, &migration.Version, &migration.MigratedOn)
            if err != nil {
                return migration, err
            }
        }
    }

    return migration, nil
}

func (dialect *MysqlDialect) CreateTable(tableName string, rawQuery string, fields MigrationFields) error {
    var err error
    if rawQuery == "" {
        rawQuery = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ("+
            "%s"+
            ")", tableName, strings.Join(fields.Ddl(), ",\n"))
    }

    _, err = dialect.connection.Query(rawQuery)
    return err
}

func (dialect *MysqlDialect) MarkDone(migration *Migration) (bool, error) {
    _, err := dialect.connection.Exec(
        fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES "+"('%s', %d, '%s')",
            "migrations",
            "migration_name",
            "migration_version",
            "migrated_on",
            migration.Name,
            migration.Version,
            time.Now().Format("2006-01-02 15:04:05"),
        ),
    )

    return err == nil, err
}