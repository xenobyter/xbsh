package view

import (
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestGetDir(t *testing.T) {

	// Prepare  test directory
	testDir, err := ioutil.TempDir("", "xbshTestDir")
	fmt.Println(testDir)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(testDir)
	testFile, err := ioutil.TempFile(testDir, "xbshTesFile.*.tst")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("getDir should return the testFile")
	want, _ := testFile.Stat()
	fmt.Println(want)
	got := getDir(testDir)[0]
	fmt.Println(got)
	if want.Name() != got.Name() {
		t.Errorf("got: %v, want: %v", got.Name(), want.Name())
	}
}
