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

type Config struct {
    DestinationPath string
}

var CommandOptions = make(map[string]interface{})

// Call with os.Args ?
func init() {
    var config Config
    // todo create Config from context or environment
    CommandOptions["init"] = func(args []string) error {
        return InitMigrations(config)
    }
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
    println("Creating migrations...")
    println(name)

    migrationCmd := strings.Split(name, "_")[0]

    switch migrationCmd {
    case "init":
    case "create":
        {
            args = append(args[1:])
            createMigrationFile(config, strings.Replace(name, "Create_", "", 9), args)

        }
    }

    return nil
}

func createMigrationFile(config Config, name string, args []string) {
    t, err := template.ParseFiles(config.DestinationPath + "/create-migration.tmpl")

    if err != nil {
        panic(err)
    }

    file, err := os.Create(fmt.Sprintf(
        "%s/%d_%s.go",
        config.DestinationPath,
        time.Now().Unix(),
        name,
    ))
    println("Filename >> ", filepath.Base(file.Name()))
    defer file.Close()

    err = t.Execute(file, nil)
    if err != nil {
        panic(err)
    }

    println("Created migration file")
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

    err := os.MkdirAll(config.DestinationPath, os.ModePerm)
    println("Creating folder", config.DestinationPath)

    return err
}

func initTemplates(config Config) error {
    _, file, _, _ := runtime.Caller(0)
    path, _ := filepath.Abs(filepath.Dir(file) + "./../../resources/templates")
    println("Copying files", path)

    err := os.CopyFS(
        config.DestinationPath,
        os.DirFS(path),
    )

    return err
}