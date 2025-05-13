package migrator

import (
    "errors"
    "fmt"
    "github.com/genellern/hela/cmd/utils"
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "text/template"
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
    action      Action
}

type MigrationInterface interface {
    Action() Action
    GetFields() []string
    TableName() string
}

func (m MigrationOptions) Action() Action {
    return m.action
}

func (m MigrationOptions) GetFields() []string {
    return m.Fields
}

func (m MigrationOptions) TableName() string {
    return m.Table
}

var CommandOptions = make(map[string]interface{})

// Call with os.Args ?
func init() {
    // todo create Config from context or environment
    //CommandOptions["init"] = func(args []string) error {
    //    return InitMigrations(config)
    //}
    CommandOptions["create"] = func(args []string) error {
        if len(args) < 2 {
            err := errors.New("Missing migration name")

            return err
        }
        return nil
    }
}

func CreateMigration(config Config, args []string) error {
    name := utils.ToSnakeCase(args[0])

    migrationCmd := strings.Split(name, "_")[0]

    switch migrationCmd {
    case "create":
        {
            return createMigrationFile(config, strings.Replace(name, "Create_", "", 9), args[1:])
        }
    }

    return nil
}

func createMigrationFile(config Config, name string, args []string) error {
    t, err := template.ParseFiles(config.DestinationPath + "/create-migration.tmpl")

    if err != nil {
        return err
    }

    file, err := os.Create(fmt.Sprintf(
        "%s/%d_%s.go",
        config.DestinationPath,
        time.Now().Unix(),
        name,
    ))
    println(filepath.Abs(file.Name()))
    defer file.Close()

    var data MigrationOptions
    data.Table = name
    data.PackageName = "migrations"
    data.Fields = args
    data.action = Create

    err = t.Execute(file, data)
    if err != nil {
        println("Couldn't create migration file >> ", filepath.Base(file.Name()))
        println(err.Error())
        return err
    }

    println("Created migration file >> ", filepath.Base(file.Name()))
    return nil
}

func InitMigrations(config Config) error {

    err := initFolders(config)
    if err != nil {
        println(err.Error())
        return err
    }

    err = initTemplates(config)
    if err != nil {
        println(err.Error())
    }

    return err
}

func initFolders(config Config) error {

    println("Creating folder", config.DestinationPath)
    return os.MkdirAll(config.DestinationPath, os.ModePerm)
}

func initTemplates(config Config) error {
    _, file, _, _ := runtime.Caller(0)
    path, _ := filepath.Abs(filepath.Dir(file) + "./../../resources/templates")
    println("Copying files", path)

    return os.CopyFS(
        config.DestinationPath,
        os.DirFS(path),
    )
}