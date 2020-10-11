package storage

import (
	"os"
	"reflect"
	"testing"
)

func TestHistoryWrite(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name   string
		cmd    string
		wantID int64
	}{
		{"should return 1 on first insert", "testing", 1},
		{"should return 2 on second insert", "testing", 2},
		{"should return 0 on empty insert", "", 0},
	}
	//setup
	dir := tempDirHelper()
	defer os.RemoveAll(dir)
	db, _ = openDB(dir + "/" + "test.sqlite")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotID := HistoryWrite(tt.cmd); gotID != tt.wantID {
				t.Errorf("HistoryWrite() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func Test_GetMaxID(t *testing.T) {
	tests := []struct {
		name   string
		wantID int64
	}{
		{"empty db, return 0", 0},
		{"should return 2", 2},
	}

	//setup
	dir := tempDirHelper()
	defer os.RemoveAll(dir)
	db, _ = openDB(dir + "/" + "test.sqlite")
	
	for i, tt := range tests {
		if i==1 {
			HistoryWrite("cmd01")
			HistoryWrite("cmd02")
		}
		t.Run(tt.name, func(t *testing.T) {
			if gotID := GetMaxID(); gotID != tt.wantID {
				t.Errorf("GetMaxID() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func TestHistoryRead(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int64
	}{
		{"should return \"\",1 when called on empty db", args{0}, "", 1},
		{"should return \"\",3 for id=3", args{3}, "", 3},
		{"should return \"cmd01\",1 for id=1", args{1}, "cmd01", 1},
		{"should return \"cmd02\",2 for id=2", args{2}, "cmd02", 2},
		{"should return \"\",3 for id=0", args{0}, "", 3},
		{"should return \"\",3 for id=-2", args{-2}, "", 3},
	}

	//setup
	dir := tempDirHelper()
	defer os.RemoveAll(dir)
	db, _ = openDB(dir + "/" + "test.sqlite")
	
	for i, tt := range tests {
		if i==1 {
			//only first test with empty db
			HistoryWrite("cmd01")
			HistoryWrite("cmd02")
		}
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := HistoryRead(tt.args.id)
			if got != tt.want {
				t.Errorf("HistoryRead() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("HistoryRead() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHistorySearch(t *testing.T) {
	tests := []struct {
		name    string
		search string
		wantRes []string
	}{
		{"should return 2 rows with search = \"\"","",[]string{"cmd02","cmd01"}},
		{"should return 2 rows with search = \"cmd0\"","cmd0",[]string{"cmd02","cmd01"}},
		{"should return 1 rows with search = \"cmd01\"","cmd01",[]string{"cmd01"}},
		{"should return 0 rows with search = \"cmd03\"","cmd03",nil},
	}

	//setup
	dir := tempDirHelper()
	defer os.RemoveAll(dir)
	db, _ = openDB(dir + "/" + "test.sqlite")
	HistoryWrite("cmd01")
	HistoryWrite("cmd02")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := HistorySearch(tt.search); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("HistorySearch() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
