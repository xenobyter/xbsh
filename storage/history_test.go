package storage

import (
	"os"
	"testing"
)

func TestHistoryWrite(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name   string
		cmd   string
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
