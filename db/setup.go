/*
Package db implements functions for local persistence of history, cache and settings
*/
package db

import (
	"database/sql"
	"log"
	"os"
	"os/user"
	"path/filepath"

	//only use init()
	_ "github.com/mattn/go-sqlite3"
	"github.com/xenobyter/xbsh/cfg"
)

var (
	db *sql.DB
)

func init() {
	dir := getDotDir()
	err := makeDotDir(dir)
	if err != nil {
		log.Panic(err)
	}
	db, err = openDB(filepath.Join(dir, cfg.DBFileName))
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dbPath string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}

	db.SetMaxOpenConns(1)
	db.Exec("PRAGMA journal_mode=WAL")

	//create inital tables
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS history (id INTEGER not null primary key, command TEXT);
		CREATE INDEX IF NOT EXISTS idx_history_command ON history (command);
		CREATE TABLE IF NOT EXISTS bin (full TEXT primary key, item TEXT, path TEXT);
		CREATE TABLE IF NOT EXISTS rename (id INTEGER not null primary key, rule TEXT);`
	_, err = db.Exec(sqlStmt)
	return
}

func getDotDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, cfg.Directory)
}

func makeDotDir(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		// dotDir doesnt exists, lets create it
		err = os.Mkdir(dir, 0755)
	}
	return err
}
