package db

import "log"

// ReadBatchRules reads all rules from database and returns them as []string
func ReadBatchRules() (lines []string) {
	var l string
	rows, err := db.Query("SELECT rule FROM batch ORDER BY id ASC;")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&l)
		if err != nil {
			log.Fatalln(err)
		}
		lines = append(lines, l)
	}
	return
}

// WriteBatchRules takes a full set of rules and writes it to db.
// Any old rules are dropped
func WriteBatchRules(lines []string) error {
	if _, err := db.Exec("DELETE FROM batch;"); err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO batch(id, rule) VALUES(?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for id, line := range lines {
		_, err = stmt.Exec(id, line)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateBatchRule takes a line number and one rule to update or insert this rule in db
func UpdateBatchRule(id int, line string) error {
	_, err := db.Exec("INSERT INTO batch(id, rule) VALUES(?, ?) ON CONFLICT(id) DO UPDATE SET rule=excluded.rule;", id, line)
	return err
}
