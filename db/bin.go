package db

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// dirScan takes a directory and writes its contents into db.
func dirScan(dir string) (err error) {
	_, err = db.Exec("DELETE FROM bin WHERE PATH = ?", dir)
	if err != nil {
		return
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		_, err = db.Exec("INSERT INTO bin(full, item, path) VALUES(?, ?, ?)",
			dir+"/"+file.Name(), file.Name(), dir)
	}
	return
}

func getPath() []string {
	return strings.Split(os.Getenv("PATH"), ":")
}

//PathCache runs throug PATH and stores files
func PathCache() {
	for _, p := range getPath() {
		dirScan(p)
	}
}

//PathComplete takes a string and returns possible completions as os.Fileinfo
func PathComplete(item string) (completions []os.FileInfo) {
	var c string
	rows, err := db.Query("SELECT full FROM bin WHERE item like ? GROUP BY item", item+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&c)
		if err != nil {
			log.Fatal(err)
		}
		f, _ := os.Stat(c)
		completions = append(completions, f)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}

//TODO: #72 Updates for bin with inotify