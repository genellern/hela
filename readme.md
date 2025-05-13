# Installation

```
go get github.com/genellern/hela
```

## Initialization
Create a config variable with the destination path for the templates and the migration files

```
config := migrator.Config{}
$config.DestinationPath = "/home/gene/GolandProjects/myapp/database/migrations"

err := migrator.InitMigrations(config)
```

## Create migrations like
```
    err := migrator.CreateMigration(config, []string{
        "CreateUsers",
        "id:uuid",
        "email:varchar[120]",
        "active:boolean",
        "created:datetime",
        "last_active:date",
        "birth_year:int?",
        "sync_time:time",
    })
```