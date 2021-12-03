package fsService

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FsService struct {
	cwd      string
	hasIndex bool
	cache    map[string][]byte
	lock     sync.Mutex
}

type FSService interface {
	TryOpen(path string)
	WriteFileContent(file *os.File, outCh chan []byte)
	WriteDirContent(file *os.File, outCh chan string)
}

func MakeFsService() (fs *FsService, err error) {

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	ext, _ := Exists(cwd + "/index.html")

	fs = &FsService{
		cwd:      cwd,
		hasIndex: ext,
		cache:    make(map[string][]byte),
	}

	// Cache the index.html directly
	if fs.hasIndex {
		dat, err := os.ReadFile("/tmp/dat")
		if err != nil {
			return fs, nil
		}

		fs.cache["index.html"] = dat
	}

	return fs, nil
}

// Possible input
// 1. abs path /var/etc/...
func (fs *FsService) TryOpen(path string) (cleanPath string, file *os.File, isDir bool, err error) {
	isDir = false

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", nil, isDir, err
	}

	res, err := filepath.Match(fs.cwd+"/*", absPath)
	if err != nil {
		return "", nil, isDir, err
	}
	if !res {
		return "", nil, isDir, errors.New("no such File")
	}

	cleanPath, err = filepath.Rel(absPath, fs.cwd)
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

func (fs *FsService) WriteFileContent(file *os.File, outCh chan []byte) (start bool, err error) {
	// Ensure the file got closed
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return false, errors.New("not a file")
	}

	defer func() {
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

	return true, nil
}

func (fs *FsService) WriteDirContent(file *os.File, outCh chan string) (start bool, err error) {
	// Ensure the file got closed
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, errors.New("not a Dir")
	}

	defer func() {
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
