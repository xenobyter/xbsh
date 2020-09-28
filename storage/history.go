package storage

import "log"

//HistoryWrite takes a command as string and stores it. HiostoryWrite returns the id for the stored command.
func HistoryWrite(cmd string) (id int64) {
	if len(cmd) == 0 {
		return
	}
	//TODO: #47 check last entry and dont store duplicates
	//TODO: #48 Trim leading spaces before storing command history
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