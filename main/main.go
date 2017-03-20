package main

import (
	"flag"
	"Zipper/zip"
)

var dirPath = flag.String("d", "./", "Provide a *relative path to the directory with the files that you want to compress.")
var suffix = flag.String("s", ".sql", "Provide a filename suffix for the algorithm.")
var prefix = flag.String("p", "", "Provide a filename prefix for the algorithm.")
var reqNum = flag.Uint("n", 0, "Required number of files to be compressed.")

func init() {
	flag.Parse()
}

func main() {
	zip.CompressFiles(*dirPath, *reqNum, *prefix, *suffix)
}
