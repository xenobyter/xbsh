package storage

import (
	"database/sql"
	"log"
)

//HistoryWrite takes a command as string and stores it. HistoryWrite returns the id for the stored command.
func HistoryWrite(cmd string) (id int64) {
	if len(cmd) == 0 {
		return
	}
	//TODO: #47 check last entry and don't store duplicates
	//TODO: #48 Trim leading spaces before storing command history
	stmtHistoryWrite, err := db.Prepare("INSERT INTO history(command) VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmtHistoryWrite.Exec(cmd)
	if err != nil {
		log.Fatal(err)
	}
	id, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return
}

//TODO: #46 Purge old entries from command history

// HistoryRead takes an id and returns the stored command and it's id.
// It returns the last inserted command and max(id) when called with id < 1
// It returns an empty string and max(id) when no command is stored for the given id
func HistoryRead(id int64) (string, int64) {
	var cmd string
	if id == 0 {
		return "", GetMaxID()+1
		// id = GetMaxID()
	}
	err := db.QueryRow("select command from history where id = ?", id).Scan(&cmd)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", GetMaxID()+1
		}
		log.Fatal(err)
	}
	return cmd, id
}

// GetMaxID returns the id for the last history item
func GetMaxID() (id int64) {
	err := db.QueryRow("select MAX(id) from history").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return
}
