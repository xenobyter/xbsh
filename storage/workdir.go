package storage

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var skipList = []string{".*/.git\\b", ".*/.cache\\b"} //TODO: #62 Add skipList to Config Management

func skipPath(path string) bool {
	for _, s := range skipList {
		match, err := regexp.Match(s, []byte(path))
		if err != nil {
			log.Fatal(err)
		}
		if match {
			return true
		}
	}
	return false
}

// WorkDirScan takes a directory and writes its contents into db.
func WorkDirScan(root string) error {
	//truncate db
	if _, err := db.Exec("DELETE FROM workdir"); err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !skipPath(path) {
			//Store found items

			if lastSlash := strings.LastIndex(path, "/"); lastSlash != -1 {
				p := path[0:lastSlash]
				i := path[lastSlash+1:]
				_, err = tx.Exec("INSERT INTO workdir(item, path, mode, isdir) VALUES(?,?,?,?)", i, p, info.Mode().Perm(), info.IsDir())
			}
		}
		return nil
	})
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return nil
}

//WorkDirCacheItem represents one cached item
type WorkDirCacheItem struct {
	item, path string
	mode       int
	isDir      bool
}

//WorkDirSearch takes a search string and returns a slice of matching WorkDirCacheItem
func WorkDirSearch(item, path string) (items []WorkDirCacheItem) {
	var wdcItem WorkDirCacheItem
	rows, err := db.Query("SELECT item, path, mode, isdir FROM workdir where item LIKE ? AND path = ?", item+"%", path)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&wdcItem.item, &wdcItem.path, &wdcItem.mode, &wdcItem.isDir)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, wdcItem)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}
