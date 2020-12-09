package db

import (
	"database/sql"
	"log"
)

//HistoryWrite takes a command as string and stores it. HistoryWrite returns the id for the stored command.
func HistoryWrite(cmd string) (id int64) {
	if lastCmd, _ := HistoryRead(GetMaxID()); len(cmd) == 0 || cmd == lastCmd {
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

// HistoryRead takes an id and returns the stored command and it's id.
// It returns an empty string and max(id) when no command is stored for the given id
func HistoryRead(id int64) (string, int64) {
	var cmd string

	err := db.QueryRow("SELECT command FROM history WHERE id = ?", id).Scan(&cmd)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", GetMaxID()
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

// CleanUp deletes old and unwanted commands from history
// delExit can be used to clean history from unwanted 'exit' commands
// maxEntires and delExit are supposed to be set from cfg
// Cleanup returns the deleted rows either from maxEntries or delExit and any error
func CleanUp(maxEntires int64, delExit string) (int64, error) {
	stmt, err := db.Prepare("DELETE FROM history WHERE command = ?;")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(delExit)
	if err != nil {
		return 0, err
	}
	r1, err := res.RowsAffected()
	if err != nil {
		return r1, err
	}

	stmt, err = db.Prepare("DELETE FROM history WHERE id IN (SELECT id FROM history ORDER BY id DESC LIMIT -1 OFFSET ?)")
	if err != nil {
		return r1, err
	}
	res, err = stmt.Exec(maxEntires)
	if err != nil {
		return r1, err
	}
	r2, err := res.RowsAffected()
	if err != nil {
		return r1 + r2, err
	}
	return r1 + r2, nil
}
