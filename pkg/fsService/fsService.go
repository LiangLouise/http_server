package fsService

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FsService struct {
	CWD      string
	HasIndex bool
	Cache    map[string][]byte
	Lock     sync.Mutex
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

	ext, _ := Exists(CWD + "/index.html")

	fs = &FsService{
		CWD:      CWD,
		HasIndex: ext,
		Cache:    make(map[string][]byte),
	}

	// Cache the index.html directly
	if fs.HasIndex {
		dat, err := os.ReadFile("/tmp/dat")
		if err != nil {
			return fs, nil
		}

		fs.Cache["index.html"] = dat
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

func (fs *FsService) WriteFileContent(file *os.File, outCh chan []byte) (start bool, size int64, err error) {

	info, err := file.Stat()
	if err != nil {
		return false, 0, err
	}

	if info.IsDir() {
		return false, 0, errors.New("not a file")
	}

	go func() {
		defer file.Close()
		reader := bufio.NewReader(file)
		// Buffer the file content into 512 bytes long buffer
		buf := make([]byte, 512)
		for {
			n, err := reader.Read(buf)

			if err != nil {

				if err != io.EOF {
					log.Printf("Error: %s", err)
				}

				break
			}

			outCh <- buf[0:n]
		}
		close(outCh)
	}()

	return true, info.Size(), nil
}

func (fs *FsService) WriteDirContent(file *os.File, outCh chan string) (start bool, err error) {

	info, err := file.Stat()
	if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, errors.New("not a Dir")
	}

	go func() {
		// Ensure the file got closed
		defer file.Close()
		files, err := file.ReadDir(-1)
		if err != nil {
			log.Printf("Error: %s", err)
			close(outCh)
			return
		}

		// Write the dir entries to output channel
		for _, file := range files {
			fileName := file.Name()
			if file.IsDir() {
				fileName += "/"
			}
			outCh <- fileName
		}
		close(outCh)
	}()

	return true, nil
}
