package fsService

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type FsService struct {
	CWD string
}

type FSService interface {
	TryOpen(path string)
	WriteFileContent(file *os.File, outCh chan []byte)
	WriteDirContent(file *os.File, outCh chan string)
}

// Create a new fs serverice instance. The working
// directory will be the same as server running dir
func MakeFsService() (fs *FsService, err error) {

	CWD, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fs = &FsService{
		CWD: CWD,
	}

	return fs, nil
}

// Possible input
// 		1. abs path /var/etc/...
// 		2. normal relatvie path a/b/c
// cases to clean
// 		1. abnormal relative path ../../a/../b/c
// 		2. abnormal relative path ../..
func (fs *FsService) TryOpen(path string) (cleanPath string, file *os.File, isDir bool, err error) {

	isDir = false
	path, err = url.QueryUnescape(path)

	if err != nil {
		return "", nil, isDir, err
	}
	if !strings.HasPrefix(path, "/") {
		return "", nil, isDir, err
	}
	// Remove the starting slash,
	// So that file system can know it's a
	// relative path, which must be under cwd
	path = strings.Replace(path, "/", "", 1)

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", nil, isDir, err
	}

	// Make sure that the parsed and cleaned
	// Abs path align with cwd, i.e.,
	// Input path: "\/\/\/\/\/" -> absPath: "\/"
	// But if absPath != $(pwd), we still need to reject
	// As the Rel will try its best to find a route match
	// target and base
	if !strings.HasPrefix(absPath, fs.CWD) {
		return "", nil, isDir, errors.New("no such File")
	}

	cleanPath, err = filepath.Rel(fs.CWD, absPath)
	if err != nil {
		return "", nil, isDir, err
	}

	f, err := os.Open(absPath)
	if err != nil {
		return "", nil, isDir, err
	}
	info, err := f.Stat()
	if err != nil {
		return "", nil, isDir, err
	}

	return cleanPath, f, info.IsDir(), nil
}

func TryOpenIndex(path string) (file *os.File, err error) {
	path = path + "/index.html"
	file, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil

}
