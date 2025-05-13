package migrator

import (
    "os"
    "path/filepath"
    "runtime"
)

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