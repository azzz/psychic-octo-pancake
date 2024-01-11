package eventslog

import (
	"fmt"
	"os"
	"sync"
	"syscall"
)

// File implements a file-base events log safe for concurrency use
// File opens a given file exclusively and locks it from other processes to guarantee
// that only the server will write consistent data.
type File struct {
	file *os.File
	sync.Mutex
}

func (f *File) Write(p []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()

	return f.file.Write(p)
}

func (f *File) Close() error {
	f.Lock()
	defer f.Unlock()

	var errs []error

	if err := funlock(f.file); err != nil {
		errs = append(errs, fmt.Errorf("unlock file: %w", err))
	}

	if err := f.file.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close file: %w", err))
	}

	f.file = nil

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// Open a file exclusively. If file does not exist then create it.
func Open(path string) (*File, error) {
	var (
		fw  = &File{}
		err error
	)

	fw.file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	if err := flock(fw.file); err != nil {
		return nil, fmt.Errorf("locking file: %w", err)
	}

	return fw, nil
}

func flock(file *os.File) error {
	flag := syscall.LOCK_NB | syscall.LOCK_EX
	if err := syscall.Flock(int(file.Fd()), flag); err != nil {
		return err
	}

	return nil
}

func funlock(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}
