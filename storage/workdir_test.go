package storage

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestWorkDirScan(t *testing.T) {
	//setup 
	dir := tempDirHelper()
	db, _ = openDB(dir + "/" + "test.sqlite")

	os.Mkdir(dir+"/empty", 0770)
	os.Mkdir(dir+"/dir1", 0770)
	os.Create(dir + "/dir1/file")
	os.MkdirAll(dir+"/dir2/dir21", 0770)
	os.Create(dir + "/dir2/dir21/file")
	os.MkdirAll(dir+"/dir3/.git", 0770)
	os.Create(dir + "/dir3/.git/file")

	tests := []struct {
		name                 string
		arg                  string
		wantErr              error
		wantItems, wantPaths []string
		wantMode             []int
		wantIsDir            []bool
	}{
		{"empty dir", dir + "/empty", nil,
			[]string{"empty"},
			[]string{dir},
			[]int{0770},
			[]bool{true}},
		{"dir with file", dir + "/dir1", nil,
			[]string{
				"dir1",
				"file"},
			[]string{
				dir,
				dir + "/dir1"},
			[]int{0770, 0664},
			[]bool{true, false}},
		{"dir with dir with file", dir + "/dir2", nil,
			[]string{
				"dir2",
				"dir21",
				"file"},
			[]string{
				dir,
				dir + "/dir2",
				dir + "/dir2/dir21"},
			[]int{0770, 0770, 0664},
			[]bool{true, true, false}},
		{"skip .git", dir + "/dir3", nil,
			[]string{"dir3"},
			[]string{dir},
			[]int{0770},
			[]bool{true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErr := WorkDirScan(tt.arg); gotErr != tt.wantErr {
				t.Errorf("WorkDirScan() = %v, want %v", gotErr, tt.wantErr)
			}
			gotItems, gotPaths, gotMode, gotIsDir := wdHelper()
			if !reflect.DeepEqual(gotItems, tt.wantItems) {
				t.Errorf("WorkDirScan(%v) stored items %v, wanted %v", tt.arg, gotItems, tt.wantItems)
			}
			if !reflect.DeepEqual(gotPaths, tt.wantPaths) {
				t.Errorf("WorkDirScan(%v) stored paths %v, wanted %v", tt.arg, gotPaths, tt.wantPaths)
			}
			if !reflect.DeepEqual(gotMode, tt.wantMode) {
				t.Errorf("WorkDirScan(%v) stored %o, wanted %o", tt.arg, gotMode, tt.wantMode)
			}
			if !reflect.DeepEqual(gotIsDir, tt.wantIsDir) {
				t.Errorf("WorkDirScan(%v) stored %v, wanted %v", tt.arg, gotIsDir, tt.wantIsDir)
			}
		})
	}

	//teardown
	os.RemoveAll(dir)
}

func wdHelper() (items, paths []string, modes []int, isDir []bool) {
	var (
		item, path string
		mode       int
		dir        bool
	)
	rows, err := db.Query("SELECT item, path, mode, isdir FROM workdir")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&item, &path, &mode, &dir)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
		paths = append(paths, path)
		modes = append(modes, mode)
		isDir = append(isDir, dir)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func TestWorkDirSearch(t *testing.T) {
	//setup
	dir := tempDirHelper()
	db, _ = openDB(dir + "/" + "test.sqlite")
	os.MkdirAll(dir+"/dir1/dir2/dir3", 0775)
	os.MkdirAll(dir+"/dir1/dir5", 0775)
	if err := WorkDirScan(dir); err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		name       string
		item, path string
		wantRes    []WorkDirCacheItem
	}{
		{"find dir1", "dir1", dir, []WorkDirCacheItem{
			{"dir1", dir, 0775, true},
		}},
		{"find dir3", "dir3", dir + "/dir1/dir2", []WorkDirCacheItem{
			{"dir3",dir + "/dir1/dir2", 0775, true},
		}},
		{"find with single wildcard", "dir", dir + "/dir1/dir2", []WorkDirCacheItem{
			{"dir3",dir + "/dir1/dir2", 0775, true},
		}},
		{"find with two wildcard", "dir", dir + "/dir1", []WorkDirCacheItem{
			{"dir2",dir + "/dir1", 0775, true},
			{"dir5",dir + "/dir1", 0775, true},
		}},
		{"dont find dir4", "dir4", dir + "/dir1/dir2", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := WorkDirSearch(tt.item,tt.path); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("WorkDirSearch(%v, %v) = %v, want %v", tt.item,tt.path, gotRes, tt.wantRes)
			}
		})
	}

	//teardown
	os.RemoveAll(dir)
}

func Test_skipPath(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{"skip .git within", "/home/.git/test", true},
		{"skip .cache last", "/home/.cache", true},
		{"skip .cache with trailing /", "/home/.cache/", true},
		{"don't skip cache", "/home/cache", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := skipPath(tt.args); got != tt.want {
				t.Errorf("skipPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
