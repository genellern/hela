package migrator

import (
    "errors"
    "os"
    "path/filepath"
    "runtime"
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