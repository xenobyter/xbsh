package db

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
		{"should return 1 on first insert", "testing1", 1},
		{"should return 2 on second insert", "testing2", 2},
		{"should return 0 on consecutive duplicate", "testing2", 0},
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
		{"2 command written, should return 2", 2},
		{"last command written again, should still return 2", 2},
	}

	//setup
	dir := tempDirHelper()
	defer os.RemoveAll(dir)
	db, _ = openDB(dir + "/" + "test.sqlite")

	for i, tt := range tests {
		switch i {
		case 1:
			HistoryWrite("cmd01")
			HistoryWrite("cmd02")
		case 2:
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
		{"should return \"\",0 when called on empty db", args{0}, "", 0},
		{"should return \"\",2 for id=3", args{3}, "", 2},
		{"should return \"cmd01\",1 for id=1", args{1}, "cmd01", 1},
		{"should return \"cmd02\",2 for id=2", args{2}, "cmd02", 2},
		{"should return \"\",2 for id=0", args{0}, "", 2},
		{"should return \"\",2 for id=-2", args{-2}, "", 2},
	}

	//setup
	dir := tempDirHelper()
	defer os.RemoveAll(dir)
	db, _ = openDB(dir + "/" + "test.sqlite")

	for i, tt := range tests {
		if i == 1 {
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
		search  string
		wantRes []string
	}{
		{"should return 2 rows with search = \"\"", "", []string{"cmd02", "cmd01"}},
		{"should return 2 rows with search = \"cmd0\"", "cmd0", []string{"cmd02", "cmd01"}},
		{"should return 1 rows with search = \"cmd01\"", "cmd01", []string{"cmd01"}},
		{"should return 0 rows with search = \"cmd03\"", "cmd03", nil},
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

func Test_cleanUp(t *testing.T) {
	type args struct {
		maxEntires int
		delExit   string
	}
	tests := []struct {
		name    string
		args    args
		wantCnt int64
		wantErr bool
	}{
		{"Delete 'exit'", args{10, "exit"}, 2, false},
		{"Delete 'delete1'", args{2, "exit"}, 1, false},
	}

	//setup
	dir := tempDirHelper()
	// defer os.RemoveAll(dir)
	db, _ = openDB(dir + "/" + "test.sqlite")
	HistoryWrite("exit")
	HistoryWrite("delete1")
	HistoryWrite("delete2")
	HistoryWrite("delete3")
	HistoryWrite("exit")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCnt, err := CleanUp(tt.args.maxEntires, tt.args.delExit)
			if (err != nil) != tt.wantErr {
				t.Errorf("cleanUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCnt != tt.wantCnt {
				t.Errorf("cleanUp() = %v, want %v", gotCnt, tt.wantCnt)
			}
		})
	}
}
