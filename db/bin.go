package db

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
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
			filepath.Join(dir, file.Name()), file.Name(), dir)
	}
	return
}

func getPath() []string {
	return strings.Split(os.Getenv("PATH"), ":")
}

//PathCache runs through PATH and stores files
func PathCache() {
	for _, p := range getPath() {
		dirScan(p)
		go makeNotify(p)
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

func makeNotify(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				processNotifyEvent(event.Name, event.Op)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func processNotifyEvent(file string, op fsnotify.Op) {
	switch {
	case op == fsnotify.Create:
		db.Exec("INSERT INTO bin(full, item, path) VALUES(?, ?, ?)",
			file,
			filepath.Base(file),
			filepath.Dir(file))
	case op == fsnotify.Remove || op == fsnotify.Rename:
		db.Exec("DELETE FROM bin WHERE full = ?", file)
	}
}
