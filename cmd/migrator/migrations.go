package migrator

import "github.com/genellern/hela/cmd/utils"

func (m *MigrationCallbacksStack) Migrate(connection *Connection) error {
    var migrations Migrations
    var err error
    latestMigration, err := connection.GetLatestMigration()
    if err != nil {
        return err
    }

    // Extract migrations
    for _, callback := range *m {
        migrations = append(migrations, callback())
    }

    migrations = migrations.filterOutDoneMigrations(latestMigration)

    // Run migrations
    for _, migration := range migrations {
        err = migration.processMigration(connection)

        if err != nil {
            return err
        }
    }

    return nil
}

func (f *MigrationFields) Ddl() []string {
    ddl := []string{}
    for _, field := range *f {
        ddl = append(ddl, field.name+" "+field.ddl)
    }

    return ddl
}

func (migrations *Migrations) filterOutDoneMigrations(latestMigration *MigrationRecord) Migrations {
    if (latestMigration) == nil {
        return *migrations
    }

    return utils.Filter(*migrations, func(migration Migration) bool {
        return migration.Version > latestMigration.Version
    })
}