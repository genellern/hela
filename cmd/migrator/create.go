package migrator

import (
    "fmt"
    "github.com/genellern/hela/cmd/utils"
    "os"
    "path/filepath"
    "strings"
    "text/template"
    "time"
)

func CreateMigration(config Config, args []string) error {
    name := utils.ToSnakeCase(args[0])

    migrationCmd := strings.Split(name, "_")[0]

    switch migrationCmd {
    case "create":
        {

            file, err := createMigrationFile(config, strings.Replace(name, "Create_", "", 9))
            if err != nil {
                return err
            }
            return parseTemplate(config, file, name, args[1:])
        }
    }

    return nil
}

func createMigrationFile(config Config, name string) (*os.File, error) {

    file, err := os.Create(fmt.Sprintf(
        time.Now().Format("2006_01_02_15_04_05")+".go",
        config.DestinationPath,
        time.Now().Unix(),
        name,
    ))
    println(filepath.Abs(file.Name()))

    return file, err
}

func parseTemplate(config Config, file *os.File, name string, args []string) error {
    t, err := template.ParseFiles(config.DestinationPath + "/create-migration.tmpl")
    defer file.Close()

    if err != nil {
        return err
    }
    var data MigrationOptions
    data.Table = name
    data.Name = string(Create) + "_" + name
    data.PackageName = "migrations"
    data.Fields = args
    data.Action = Create

    err = t.Execute(file, data)
    if err != nil {
        println("Couldn't create migration file >> ", filepath.Base(file.Name()))
        println(err.Error())
        return err
    }

    println("Created migration file >> ", filepath.Base(file.Name()))
    return nil
}