package zip

import (
	"os"
	"regexp"
	"bytes"
	"archive/zip"
	"io/ioutil"
	"log"
	"path/filepath"
)

//A FileCollection struct that eases the operations, but mostly visually. :)
type FileCollection struct {
	Files           []os.FileInfo
	CompressedFiles []os.FileInfo
}

//Filter a FileCollection by the given regex
func filter(regexString string, fc FileCollection) FileCollection {
	var fCollection = new(FileCollection)

	r, err := regexp.Compile(regexString)

	if err != nil {
		log.Fatalf("Failed to compile regex: %s. Err: %e", regexString,  err)
	}

	for _, file := range fc.Files {
		if r.MatchString(file.Name()) {
			fCollection.Files = append(fCollection.Files, file)
		}
	}

	return *fCollection
}

func (f FileCollection) filterPrefix(prefix string) FileCollection {
	return filter("^(" + prefix +")", f)
}

func (f FileCollection) filterSuffix(suffix string) FileCollection {
	return filter(suffix + "$", f)
}

func CompressFiles(dirPath string, reqNum uint, prefix, suffix string) {

	buf := new(bytes.Buffer)

	w := zip.NewWriter(buf)

	files, err := ioutil.ReadDir(dirPath)

	if err != nil {
		log.Fatal("Couldn't read files from directory. Err: ", err)
	}

	//A FileCollection that we'll work on
	fCollection := FileCollection{files, nil}

	//Apply filters to it
	fCollection = fCollection.filterPrefix(prefix).filterSuffix(suffix)

	if filesLength := len(fCollection.Files); filesLength < int(reqNum) {
		log.Fatalf("Not enough files to be zipped. Found %d, min: %d\n", filesLength, reqNum)
	}

	for i := 0; i < int(reqNum); i++ {

		// The name must be a relative path: it must not start with a drive
		// letter (e.g. C:) or leading slash, and only forward slashes are
		// allowed.
		f, err := w.Create(fCollection.Files[i].Name())

		if err != nil {
			log.Fatal("Error creating a writer. Err: ", err)
		}

		fData, err := ioutil.ReadFile(filepath.Join(dirPath, fCollection.Files[i].Name()))

		if err != nil {
			log.Fatal("Error reading a file. Err: ", err)
		}

		_, err = f.Write(fData)

		if err != nil {
			log.Fatal("Error writing to a file. Err: ", err)
		}

		//Append to the slice of compressed files
		fCollection.CompressedFiles = append(fCollection.CompressedFiles, fCollection.Files[i])
	}

	if err := w.Close(); err != nil {
		log.Fatal("Error closing writer! Err: ", err)
	}

	//Take the first file for the first part of the zipped file name
	fileName := fCollection.CompressedFiles[0].Name()

	cFilesLength := len(fCollection.CompressedFiles)

	//However, if there's more than 1 file, append to it the name of the last file
	if cFilesLength > 1 {
		fileName += "-" + fCollection.CompressedFiles[cFilesLength-1].Name()
	}

	if err := ioutil.WriteFile(filepath.Join(dirPath, fileName + ".zip"), buf.Bytes(), 066); err != nil {
		log.Fatalf("Error writing to a file. Filename: %s, Err: %s", fileName, err)
	}

	for _, file := range fCollection.CompressedFiles {
		if err := os.Remove(filepath.Join(dirPath, file.Name())); err != nil {
			log.Printf("Failed to remove a file. Filename: %s, Err: %s", file.Name(), err)
		}
	}
}