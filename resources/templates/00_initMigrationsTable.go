package migrations

import "github.com/genellern/hela/cmd/migrator"

func (m *Migration) Migration_0000_InitMigrationsTable() migrator.Migration {
    return migrator.Migration{
        Version:     0,
        Name:        "Init migrations",
        Table:       "migrations",
        PackageName: "migrations",
        Action:      migrator.Create,
        Fields: []string{
            "migration_name:varchar(120)",
            "migration_version:int(5)",
            "migrated_on:datetime",
        },
    }
}