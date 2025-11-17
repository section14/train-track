package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Env struct {
    Db *sql.DB
}

var dbLocation = "/Users/anthony/embedded/train-track/.data/db.sqlite"

func NewEnv() *Env {
    env := &Env{}

    db, err := sql.Open("sqlite3", dbLocation)
    if err != nil {
        log.Fatal("couldn't connect to db", err)
    }

    env.Db = db

    return env
}
