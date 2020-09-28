package storage

import (
	"log"
	"database/sql"
	"os"
	"os/user"

	//only use init()
	_ "github.com/mattn/go-sqlite3"
)

const (
	dotDirSuffix = ".xbsh"
	dbFileName   = "storage.sqlite"
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
	db, err = openDB(dir + "/" + dbFileName)
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dbPath string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}

	//create inital tables
	sqlStmt := "CREATE TABLE IF NOT EXISTS history (id INTEGER not null primary key, command TEXT);"
	_, err = db.Exec(sqlStmt)
	return
}

func getDotDir() string {
	usr, _ := user.Current()
	return usr.HomeDir + "/" + dotDirSuffix
}

func makeDotDir(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		// dotDir doesnt exists, lets create it
		err = os.Mkdir(dir, 0755)
	}
	return err
}
