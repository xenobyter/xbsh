package exec

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_job_start(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		args       []string
		wantErr    bool
		wantStdOut string
		wantStdErr string
	}{
		{"Simple echo", 0, []string{"echo", "test"}, false, "test\n", ""},
		{"Wrong pwd", 1, []string{"pwd", "-v"}, true, "", "pwd: invalid option -- 'v'\nTry 'pwd --help' for more information.\n"},
		{"No cmd", 2, []string{""}, true, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := newBgJob(tt.args)
			err := j.start()

			fmt.Println(string(j.stderr))
			if (err != nil) != tt.wantErr {
				t.Errorf("job.start() error = %v, wantErr %v", err, tt.wantErr)
			}
			if string(j.stdout) != tt.wantStdOut {
				t.Errorf("job.start() stdout = %v, wantStdOut %v", string(j.stdout), tt.wantStdOut)
			}
			if string(j.stderr) != tt.wantStdErr {
				t.Errorf("job.start() stderr = %v, wantStdErr %v", string(j.stderr), tt.wantStdErr)
			}
		})
	}
}

func TestBgJob(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantID   int
		wantCmd  string
		wantArgs []string
	}{
		{"First job", []string{"job1"}, 0, "job1", []string{}},
		{"Second job with args", []string{"job2", "arg0", "arg1"}, 1, "job2", []string{"arg0", "arg1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotID := BgJob(tt.args); gotID != tt.wantID {
				t.Errorf("BgJob() = %v, want %v", gotID, tt.wantID)
			}
			if gotCmd := jobList[tt.wantID].cmd; gotCmd != tt.wantCmd {
				t.Errorf("BgJob() = %v, want %v", gotCmd, tt.wantCmd)
			}
			if gotArgs := jobList[tt.wantID].args; !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("BgJob() = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestBgList(t *testing.T) {
	tests := []struct {
		name     string
		job      []string
		wantCmds []string
	}{
		{"Simple ls", []string{"ls"}, []string{"ls"}},
		{"Second cmd with args", []string{"cmd", "arg"}, []string{"ls", "cmd"}},
	}

	//Setup
	jobList = nil

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BgJob(tt.job)
			gotCmds := BgList()
			if !reflect.DeepEqual(gotCmds, tt.wantCmds) {
				t.Errorf("BgList() gotCmds = %v, want %v", gotCmds, tt.wantCmds)
			}
		})
	}
}

func TestBgGet(t *testing.T) {
	tests := []struct {
		name         string
		id           int
		wantCmd      string
		wantArgs     []string
		wantStdout   []byte
		wantStderr   []byte
		wantFinished bool
	}{
		{"Wrong id", 2, "", nil, nil, nil, false},
		{"Simple echo", 0, "echo", []string{"test"}, []byte("test\n"), []byte{}, true},
		{"Wrong args", 1, "pwd", []string{"-v"}, []byte{}, []byte("pwd: invalid option -- 'v'\nTry 'pwd --help' for more information.\n"), true},
	}

	//setup
	jobList = nil
	BgJob([]string{"echo", "test"})
	BgJob([]string{"pwd", "-v"})
	time.Sleep(time.Second)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotArgs, gotStdout, gotStderr, gotFinished := BgGet(tt.id)
			if gotCmd != tt.wantCmd {
				t.Errorf("BgGet() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("BgGet() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
			if !reflect.DeepEqual(gotStdout, tt.wantStdout) {
				t.Errorf("BgGet() gotStdout = %v, want %v", gotStdout, tt.wantStdout)
			}
			if !reflect.DeepEqual(gotStderr, tt.wantStderr) {
				t.Errorf("BgGet() gotStderr = %v, want %v", gotStderr, tt.wantStderr)
			}
			if gotFinished != tt.wantFinished {
				t.Errorf("BgGet() gotFinished = %v, want %v", gotFinished, tt.wantFinished)
			}
		})
	}
}

func TestBgDelete(t *testing.T) {
	jobList=nil
	BgJob([]string{"echo", "test"})
	BgJob([]string{"echo", "test"})
	BgJob([]string{"echo", "test"})
	time.Sleep(time.Second)
	tests := []struct {
		name string
		id int
		wantList []*job
	}{
		{"delete middle", 1, []*job{jobList[0], jobList[2]}},
		{"delete last", 1, []*job{jobList[0]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BgDelete(tt.id)
			if !reflect.DeepEqual(jobList, tt.wantList) {
				t.Errorf("BgDelete() gotList = %v, want %v", jobList, tt.wantList)
			}
		})
	}
}
