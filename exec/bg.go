package exec

import (
	"io/ioutil"
	"os"
	"os/exec"
)

// job holds a single background job
type job struct {
	process        *os.Process
	cmd            string
	args           []string
	stderr, stdout []byte
	finished       bool
}

var (
	jobList []*job
)

// BgJob starts a new job and returns it's id
func BgJob(args []string) int {
	j := newBgJob(args)
	jobList = append(jobList, j)
	go j.start()
	return len(jobList) - 1
}

// BgList returns a list of all known jobs
func BgList() (cmds []string) {
	for _, j := range jobList {
		cmds = append(cmds, j.cmd)
	}
	return
}

// BgGet takes the id of a job and returns the content for this job
func BgGet(id int) (cmd string, args []string, stdout, stderr []byte, finished bool) {
	if id >= len(jobList) {
		return
	}
	j := jobList[id]
	return j.cmd, j.args, j.stdout, j.stderr, j.finished
}

func newBgJob(args []string) *job {
	return &job{cmd: args[0], args: args[1:]}
}

// BgDelete takes an id and deletes the coressponding job from jobList.
// When the given job is still running, it will be killed.
func BgDelete(id int) {
	j:=jobList[id]
	if !j.finished {
		j.kill()
	}
	jobList=remove(jobList,id)
}

func remove(slice []*job, id int) []*job {
    return append(slice[:id], slice[id+1:]...)
}

func (j *job) start() (err error) {
	cmd := exec.Command(j.cmd, j.args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	j.process = cmd.Process
	j.stdout, err = ioutil.ReadAll(stdout)
	j.stderr, err = ioutil.ReadAll(stderr)

	err = cmd.Wait()

	j.finished = true
	return
}

func (j *job) kill() {
		j.process.Kill()
}
