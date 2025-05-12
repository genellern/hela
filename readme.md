# Installation

```
go get github.com/genellern/hela
```

## Initialization
Create a config variable with the destination path for the templates and the migration files

```
config := migrator.Config{}
$config.DestinationPath = "/home/gene/GolandProjects/permitto/database/migrations"

err := migrator.InitMigrations(config)
```

## Create migrations like
```
err = migrator.CreateMigration(config, []string{"CreateAcos"})
```