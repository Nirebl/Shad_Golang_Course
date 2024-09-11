//go:build !solution

package fileleak

import (
	"io/fs"
	"os"
	"reflect"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func getCurrFiles() []fs.FileInfo {
	filesDirEntry, err := os.ReadDir("/proc/self/fd")

	if err != nil {
		panic(err.Error())
	}

	files := make([]fs.FileInfo, 0)
	for i := range filesDirEntry {
		info, err := filesDirEntry[i].Info()

		if err != nil {
			continue
		}

		files = append(files, info)
	}

	return files
}

func VerifyNone(t testingT) {
	start := getCurrFiles()
	t.Cleanup(func() {
		end := getCurrFiles()
		if !reflect.DeepEqual(end, start) {
			t.Errorf("ERROR")
		}
	})
}
