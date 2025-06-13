# Installation

```
go get github.com/genellern/hela
```

## Initialization
Create a config variable with the destination path for the templates and the migration files

```
config := migrator.Config{}
connection := &migrator.Connection{
    Driver: migrator.MySQL,
    DSN: "developer:secret@tcp(127.0.0.1:3306)/mydb?parseTime=true",
}

$config.DestinationPath = "/home/gene/GolandProjects/myapp/database/migrations"

err := migrator.InitMigrations(config)
```

## Create migrations like
```
    err := migrator.CreateMigration(config, []string{
        "CreateUsers",
        "id:uuid",
        "email:varchar(120)|unique",
        "active:boolean",
        "created:datetime",
        "last_active:datetime|nullable",
        "birth_year:int(11)|nullable",
        "sync_time:time",
    })
```
It will create a file with a migration that looks somewhat like this:
```
package migrations

import (
    "github.com/genellern/hela/cmd/migrator"
)

func (m *Migration) Migration1747340138CreateUsers() migrator.Migration {

    return migrator.Migration{
        Version:     2,
        Name:        "Create Users",
        Table:       "users",
        PackageName: "migrations",
        Action:      migrator.Create,
        Fields: []string{
            "id:uuid",
            "username:text|unique",
            "password:text|nullable",
            "active:boolean",
            "created:datetime",
            "last_active:datetime|nullable",
            "birth_year:int(11)|nullable",
            "sync_time:time",
        },
    }
}
```

Change this file as needed, the creation call can be removed afterwards, you don't need to create them over and over.
This is more of an utilitary call.

## Add migrations to the stack
```
var migrationsStack migrator.MigrationCallbacksStack

func addMigrationCallback(migrationCallback migrator.MigrationCallback) {
    migrationsStack = append(migrationsStack, migrationCallback)
}

Then...
addMigrationCallback(m.Migration_0000_InitMigrationsTable)
addMigrationCallback(m.Migration1747340138CreateUsers)
addMigrationCallback(m.Migration_1747340138CreatePosts)
```

## Run the migrations
```
fmt.println("Migrating...")

err := migrationsStack.Migrate(connection)
if err != nil {
    panic(err)
}
defer func(Conn migrator.Connection) {
    err := Conn.Close()
    if err != nil {
        panic(err)
    }
}(config.Conn)
```