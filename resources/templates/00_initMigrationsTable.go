package migrations

import "github.com/genellern/hela/cmd/migrator"

func (m *Migration) Migration_0000_InitMigrationsTable() migrator.MigrationOptions {
    return migrator.MigrationOptions{
        Version:     0,
        Table:       "migrations",
        PackageName: "migrations",
        Action:      migrator.Create,
        Fields: []string{
            "migration_name:string",
            "migration_version:int(5)",
            "migrated_on:datetime",
        },
    }
}