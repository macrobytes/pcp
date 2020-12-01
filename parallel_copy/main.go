package main

import (
	"flag"
	"os"

	"../fileutil"
)

func printUsage() {
	println("usage: cp --src SOURCE --dest DEST [--threads N]")
}

func main() {
	threads := flag.Int("threads", 1,
		"specify the number of threads for copy operation")
	src := flag.String("src", "", "source directory")
	dest := flag.String("dest", "", "destination directory")
	flag.Parse()

	if len(*src) == 0 || len(*dest) == 0 {
		printUsage()
		os.Exit(1)
	}

	fileutil.ParallelCopy(*src, *dest, *threads)
}
