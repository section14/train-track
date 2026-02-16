package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Env struct {
    Db *sql.DB
}

func NewEnv() *Env {
    env := &Env{}
    wd, _ := os.Getwd()

    dbLocation := fmt.Sprintf("%s/.data/db.sqlite", wd)

    db, err := sql.Open("sqlite3", dbLocation)
    if err != nil {
        log.Fatal("couldn't connect to db", err)
    }

    env.Db = db

    return env
}
