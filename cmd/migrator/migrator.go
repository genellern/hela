package migrator

type Action string

const (
    Create Action = "create"
    Alter  Action = "update"
    Drop   Action = "drop"
)

type Dialect string

const (
    MySQL Dialect = "MySQL"
)

type Config struct {
    DestinationPath string
    Dialect         Dialect
}

type MigrationOptions struct {
    Table       string
    PackageName string
    Fields      []string
    Action      Action
}

type MigrationInterface interface {
    GetFields() []string
    TableName() string
}

func (m MigrationOptions) GetFields() []string {
    return m.Fields
}

func (m MigrationOptions) TableName() string {
    return m.Table
}