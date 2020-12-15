package db

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadRenameRules(t *testing.T) {
	dir, err := ioutil.TempDir("", "xbsh")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(dir)
	db, _ = openDB(filepath.Join(dir, "test.sqlite"))
	tests := []struct {
		name      string
		wantLines []string
		setLine   string
	}{
		{"Empty slice on fresh start", nil, ""},
		{"Insert first rule", []string{"rule1"}, "rule1"},
		{"Insert second rule", []string{"rule1", "rule2"}, "rule2"},
	}
	for _, tt := range tests {
		if tt.setLine != "" {
			if _, err := db.Exec("INSERT INTO batch(rule) VALUES (?);", tt.setLine); err != nil {
				log.Fatalln(err)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			if gotLines := ReadBatchRules(); !reflect.DeepEqual(gotLines, tt.wantLines) {
				t.Errorf("ReadRenameRules() = %v, want %v", gotLines, tt.wantLines)
			}
		})
	}
}

func TestWriteRenameRules(t *testing.T) {
	dir, err := ioutil.TempDir("", "xbsh")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(dir)
	db, _ = openDB(filepath.Join(dir, "test.sqlite"))
	tests := []struct {
		name    string
		lines   []string
		wantErr bool
		oldRule string
	}{
		{"Overwrite all existing rules", []string{"rule1", "rule2"}, false, "rule0"},
		{"Writing empty rules", nil, false, "rule0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := db.Exec("INSERT INTO batch(rule) VALUES (?);", tt.oldRule); err != nil {
				log.Fatalln(err)
			}
			if err := WriteBatchRules(tt.lines); (err != nil) != tt.wantErr {
				t.Errorf("WriteRenameRules() error = %v, wantErr %v", err, tt.wantErr)
			}
			rules := ReadBatchRules()
			if !reflect.DeepEqual(tt.lines, rules) {
				t.Errorf("WriteRenameRules() rules = %v, want %v", rules, tt.lines)
			}
		})
	}
}

func TestUpdateRenameRule(t *testing.T) {
	dir, err := ioutil.TempDir("", "xbsh")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(dir)
	db, _ = openDB(filepath.Join(dir, "test.sqlite"))

	type args struct {
		id   int
		line string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantRules []string
	}{
		{"Insert first rule", args{0, "rule0"}, false, []string{"rule0"}},
		{"Insert second rule", args{1, "rule1"}, false, []string{"rule0", "rule1"}},
		{"Update second rule", args{1, "rule1update"}, false, []string{"rule0", "rule1update"}},
		{"Insert empty rule", args{2, ""}, false, []string{"rule0", "rule1update", ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateBatchRule(tt.args.id, tt.args.line); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRenameRule() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotRules := ReadBatchRules()
			if !reflect.DeepEqual(tt.wantRules, gotRules) {
				t.Errorf("UpdateRenameRule() rules = %v, wantRules %v", gotRules, tt.wantRules)
			}
		})
	}
}
