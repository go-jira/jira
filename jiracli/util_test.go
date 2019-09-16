package jiracli

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func realPath(path string) string {
	cpath, err := filepath.EvalSymlinks(path)
	if err != nil {
		log.Fatal(err)
	}
	return cpath
}

func comparePaths(p1 string, p2 string) bool {
	if realPath(p1) == realPath(p2) {
		return true
	}
	return false
}

func TestFindClosestParentPath(t *testing.T) {
	dir, err := ioutil.TempDir("", "testFindParentPath")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	origDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	defer os.Chdir(origDir)

	t1 := filepath.Join(dir, "/.test1")
	t2 := filepath.Join(t1, "/.test2")
	t3 := filepath.Join(t2, "/.test1")
	err = os.MkdirAll(t3, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(t3)

	path1, err := findClosestParentPath(".test1")
	if err != nil {
		t.Errorf("findClosestParentPath should not have errored: %s", err)
	}
	if ok := comparePaths(path1, t3); !ok {
		t.Errorf("%s != %s", path1, t3)
	}

	path2, err := findClosestParentPath(".test2")
	if err != nil {
		t.Errorf("findClosestParentPath should not have errored: %s", err)
	}
	if ok := comparePaths(path2, t2); !ok {
		t.Errorf("%s != %s", path2, t2)
	}

	path3, err := findClosestParentPath(".test3")
	if err.Error() != ".test3 not found in parent directory hierarchy" {
		t.Errorf("incorrect error from findClosestParentPath: %s", err)
	}
	if path3 != "" {
		t.Errorf("path3 should be empty, but is not: %s", path3)
	}

}
