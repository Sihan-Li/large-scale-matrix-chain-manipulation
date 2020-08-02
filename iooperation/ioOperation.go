package iooperation

import (
	"fmt"
	"os"
	"io/ioutil"
	"path"
)

//CreateFile from path
func CreateFile(path string) {
	// detect if file exists
	var _, err = os.Stat(path)
	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if IsError(err) {
			return
		}
		defer file.Close()
	}
}

//IsError reports the error
func IsError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

//WriteFile write the data into file*/
func WriteFile(p [][]float64, path string) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if IsError(err) {
		return
	}
	defer file.Close()
	// write into file
	_, err = file.WriteString(fmt.Sprintln(p))
	if IsError(err) {
		return
	}
	// save changes
	err = file.Sync()
	if IsError(err) {
		return
	}
}

//ClearDir is to delete all files
func ClearDir(dir string) error {
	names, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entery := range names {
		os.RemoveAll(path.Join([]string{dir, entery.Name()}...))
	}
	return nil
}
