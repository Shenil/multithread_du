package multithread_du_test

import "testing"
import "os"
import "multithread_du"

func AssertEqual(t *testing.T, actualValue int64, expectedValue int64) {
	if actualValue != expectedValue {
		t.Fatalf("Expected %d, got %d", expectedValue, actualValue)
	}
}

func MkdirP(path string) {
	err := os.MkdirAll(path, os.FileMode(int(0777)))
	if err != nil {
		panic("Error creating directory")
	}
}

func ClearDir(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		panic("Error clearing temp directory")
	}
}

func MakeFile(path string, size int) {
	data := make([]byte, size)

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	n, err := f.Write(data)
	if err != nil {
		panic(err)
	}
	if n != size {
		panic("Didn't write all bytes to file")
	}
}

func MakeTempDir() string {
	path := "tmp/multithread_du_test"
	ClearDir(path)
	MkdirP(path)
	return path
}

func TestEmptyDir(t *testing.T) {
	dirname := MakeTempDir()

	x := multithread_du.TotalFileSize(dirname)
	AssertEqual(t, x, 0)
}

func TestSimple(t *testing.T) {
	const blocksize = 512
	dirname := MakeTempDir()
	MkdirP(dirname + "/a/1")
	MkdirP(dirname + "/b")
	MakeFile(dirname+"/a/1/file1", blocksize*11)
	MakeFile(dirname+"/a/1/file2", blocksize*20)
	MakeFile(dirname+"/file3", blocksize*15)

	x := multithread_du.TotalFileSize(dirname)
	AssertEqual(t, x, 56)
}

// don't double count hard links
func TestHardLinks(t *testing.T) {
	const blocksize = 512
	dirname := MakeTempDir()
	MkdirP(dirname + "/a/1")
	MkdirP(dirname + "/b")
	MakeFile(dirname+"/a/1/file1", blocksize*11)
	MakeFile(dirname+"/a/1/file2", blocksize*20)
	MakeFile(dirname+"/file3", blocksize*15)
	err := os.Link(dirname+"/file3", dirname+"/hardlink")
	if err != nil {
		panic(err)
	}

	x := multithread_du.TotalFileSize(dirname)
	AssertEqual(t, x, 56)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
