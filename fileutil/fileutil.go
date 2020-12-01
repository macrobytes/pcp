package fileutil

import (
	"container/list"
	"io"
	"log"
	"os"
	"path/filepath"

	"../threadpool"
)

type fileStat struct {
	path     string
	fileInfo os.FileInfo
}

type makeDirRunner struct {
	path       string
	permission os.FileMode
}

func (e *makeDirRunner) Run() {
	err := os.MkdirAll(e.path, e.permission)
	if err != nil {
		log.Fatal(err)
	}
}

func createMakeDirRunner(path string, permission os.FileMode) *makeDirRunner {
	var runner makeDirRunner
	runner.path = path
	runner.permission = permission.Perm()
	return &runner
}

type copyFileRunner struct {
	src        string
	dest       string
	permission os.FileMode
}

func createFileCopyRunner(src string, dest string, permission os.FileMode) *copyFileRunner {
	var runner copyFileRunner
	runner.src = src
	runner.dest = dest
	runner.permission = permission
	return &runner
}

func (e *copyFileRunner) Run() {
	srcReader, err := os.Open(e.src)
	if err != nil {
		log.Fatal(err)
	}
	defer srcReader.Close()

	destWriter, err := os.Create(e.dest)
	if err != nil {
		log.Fatal(err)
	}
	defer destWriter.Close()

	_, err = io.Copy(destWriter, srcReader)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chmod(e.dest, e.permission)
	if err != nil {
		log.Fatal(err)
	}
}

// ParallelCopy - Leverages the thread pool executor to
// copy the source directory to the specified destination
func ParallelCopy(src string, dest string, threads int) {
	executor := threadpool.Init(threads)

	files, dirs := readDir(src)
	for dir := dirs.Front(); dir != nil; dir = dir.Next() {
		stat := dir.Value.(fileStat)
		path := dest + "/" + stat.path[len(src):]
		runner := createMakeDirRunner(path, stat.fileInfo.Mode())
		executor.Execute(runner)
	}
	executor.Wait()

	for file := files.Front(); file != nil; file = file.Next() {
		stat := file.Value.(fileStat)
		dest := dest + "/" + stat.path[len(src):]
		runner := createFileCopyRunner(stat.path, dest, stat.fileInfo.Mode())
		executor.Execute(runner)
	}
	executor.Wait()
}

func readDir(path string) (*list.List, *list.List) {
	fileList := list.New()
	dirList := list.New()

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirList.PushBack(fileStat{
				path:     path,
				fileInfo: info,
			})
		} else {
			fileList.PushBack(fileStat{
				path:     path,
				fileInfo: info,
			})
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return fileList, dirList
}
