package db

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func binHelper() (items, paths []string) {
	var (
		item, path string
	)
	rows, err := db.Query("SELECT item, path FROM bin")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&item, &path)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
		paths = append(paths, path)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Test_dirScan(t *testing.T) {
	//setup
	dir := tempDirHelper()
	defer os.RemoveAll(dir)
	db, _ = openDB(dir + "/" + "test.sqlite")
	os.Mkdir(dir+"/bin", os.ModePerm)
	os.Create(dir + "/bin/cmd")

	tests := []struct {
		name                 string
		dir                  string
		wantErr              bool
		wantItems, wantPaths []string
	}{
		{"Cmd once", dir + "/bin", false, []string{"cmd"}, []string{dir + "/bin"}},
		{"Cmd again", dir + "/bin", false, []string{"cmd"}, []string{dir + "/bin"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := dirScan(tt.dir); (err != nil) != tt.wantErr {
				t.Errorf("dirScan() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotItems, gotPaths := binHelper()
			if !reflect.DeepEqual(tt.wantItems, gotItems) {
				t.Errorf("dirScan() items = %v, wantItems %v", gotItems, tt.wantItems)
			}
			if !reflect.DeepEqual(tt.wantPaths, gotPaths) {
				t.Errorf("dirScan() items = %v, wantPaths %v", gotPaths, tt.wantPaths)
			}
		})
	}
}

func Test_getPath(t *testing.T) {
	//Setup
	p := os.Getenv("PATH")
	os.Setenv("PATH", "test1:test2")
	defer os.Setenv("PATH", p)

	tests := []struct {
		name  string
		wantP []string
	}{
		{"Test Path", []string{"test1", "test2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotP := getPath(); !reflect.DeepEqual(gotP, tt.wantP) {
				t.Errorf("getPath() = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}

func TestPathComplete(t *testing.T) {
	//setup
	dir := tempDirHelper()
	db, _ = openDB(dir + "/" + "test.sqlite")
	defer os.RemoveAll(dir)
	os.Mkdir(dir+"/bin", os.ModePerm)
	f1, _ := os.Create(dir + "/bin/file1")
	fI1,_:= f1.Stat()
	f2, _ := os.Create(dir + "/bin/file2")
	fI2,_:=f2.Stat()
	dirScan(dir + "/bin")
	type args struct {
		item string
	}
	tests := []struct {
		name            string
		args            args
		wantCompletions []os.FileInfo
	}{
		{"Return nil when no match", args{"file3"}, nil},
		{"Return f1,f2 on empty arg", args{""}, []os.FileInfo{fI1,fI2}},
		{"Return f1", args{"file1"}, []os.FileInfo{fI1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCompletions := PathComplete(tt.args.item); !reflect.DeepEqual(gotCompletions, tt.wantCompletions) {
				t.Errorf("PathComplete() = %v, want %v", gotCompletions, tt.wantCompletions)
			}
		})
	}
}