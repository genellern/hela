package migrator

import (
    "fmt"
    "github.com/genellern/hela/cmd/utils"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "text/template"
    "time"
)

func CreateMigration(config Config, args []string) error {
    name := utils.ToSnakeCase(args[0])

    migrationCmd := strings.Split(name, "_")[0]
    timestamp := time.Now().Format("20060102_1504")

    switch migrationCmd {
    case "create":
        {

            file, err := createMigrationFile(config, timestamp+"_"+strings.Replace(name, "Create_", "", 9))
            if err != nil {
                return err
            }
            return parseTemplate(config, file, timestamp+"_"+name, args[1:])
        }
    }

    return nil
}

func createMigrationFile(config Config, name string) (*os.File, error) {

    file, err := os.Create(fmt.Sprintf(
        "%s/%s.go",
        config.DestinationPath,
        name,
    ))

    return file, err
}

func parseTemplate(config Config, file *os.File, name string, args []string) error {
    t, err := template.ParseFiles(config.DestinationPath + "/create-migration.tmpl")
    defer file.Close()

    if err != nil {
        return err
    }
    var data Migration
    data.Version, _ = strconv.Atoi(time.Now().Format("0102150405"))
    data.Table = name
    data.Name = name
    data.PackageName = "migrations"
    data.Fields = args
    // TODO get as argument
    data.Action = Create

    err = t.Execute(file, data)
    if err != nil {
        println("Couldn't create migration file >> ", filepath.Base(file.Name()))
        return err
    }

    println("Created migration file >> ", filepath.Base(file.Name()))
    return nil
}