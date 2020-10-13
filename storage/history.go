package storage

import (
	"database/sql"
	"log"
)

//HistoryWrite takes a command as string and stores it. HistoryWrite returns the id for the stored command.
func HistoryWrite(cmd string) (id int64) {
	if lastCmd, _ := HistoryRead(-1); len(cmd) == 0 || cmd == lastCmd {
		return
	}
	stmt, err := db.Prepare("INSERT INTO history(command) VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(cmd)
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
// It returns an empty string and max(id)+1 when no command is stored for the given id
func HistoryRead(id int64) (string, int64) {
	var cmd string
	if id == -1 {
		id = GetMaxID()
	}
	err := db.QueryRow("SELECT command FROM history WHERE id = ?", id).Scan(&cmd)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", GetMaxID() + 1
		}
		log.Fatal(err)
	}
	return cmd, id
}

// GetMaxID returns the id for the last history item
func GetMaxID() (id int64) {
	err := db.QueryRow("SELECT MAX(id) FROM history").Scan(&id)
	if err != nil {
		return 0
	}
	return
}

// HistorySearch takes a string and returns matching commands from history as slice of strings
func HistorySearch(search string) (res []string) {
	var cmd string
	rows, err := db.Query("SELECT command FROM history WHERE command LIKE ? ORDER BY id DESC", "%"+search+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&cmd)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, cmd)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}
