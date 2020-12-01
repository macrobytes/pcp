package fileutil

import (
	"os"
	"testing"
)

func TestGetFiles(t *testing.T) {
	src := "../sample_data"
	files, dirs := readDir(src)
	if files.Len() != 5 {
		t.Errorf("expected 5 files, got %v", files.Len())
	}

	if dirs.Len() != 5 {
		t.Errorf("expected 5 directories, got %v", dirs.Len())
	}
}

func TestParallelCopy(t *testing.T) {
	src := "../sample_data"
	dest := "/tmp/parallel-copy-test"
	testDirPath := "/tmp/parallel-copy-test/a/b/c/d"
	testFilePath := "/tmp/parallel-copy-test/a/b/c/d/4"
	testEmptyFilePath := "/tmp/parallel-copy-test/a/b/c/d/5"

	os.RemoveAll(dest)

	ParallelCopy(src, dest, 5)
	info, err := os.Stat(testDirPath)
	if err != nil {
		t.Error(err)
	}
	if !info.IsDir() {
		t.Errorf("failed to generate test directory path: %v", testDirPath)
	}

	info, err = os.Stat(testFilePath)
	if err != nil {
		t.Error(err)
	}
	if info.IsDir() {
		t.Errorf("failed to generate test file path: %v", testFilePath)
	}
	if info.Size() != 2 {
		t.Errorf("expected file-size 2, got %v", info.Size())
	}

	info, err = os.Stat(testEmptyFilePath)
	if err != nil {
		t.Error(err)
	}
	if info.IsDir() {
		t.Errorf("failed to generate test file path: %v", testEmptyFilePath)
	}
	if info.Size() != 0 {
		t.Errorf("expected file-size 0, got %v", info.Size())
	}
	os.RemoveAll(dest)
}
