package tidy

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// File represents a file on disk.
type File struct {
	Name        string
	Path        string
	IsDirectory bool
	Extension   string
	Age         time.Duration

	sizeLock  sync.Mutex
	SizeBytes int64
}

// NewFile initializes and returns a new file.
func NewFile(path string, info os.FileInfo) (*File, error) {
	f := File{
		Path:        path,
		IsDirectory: info.IsDir(),
		Name:        strings.TrimSuffix(info.Name(), filepath.Ext(path)),
		Extension:   strings.TrimPrefix(filepath.Ext(path), "."),
		Age:         time.Since(info.ModTime()),
	}

	if f.IsDirectory {
		err := f.directorySize()
		if err != nil {
			return nil, err
		}
	} else {
		f.SizeBytes = info.Size()
	}
	return &f, nil
}

func (f *File) walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		f.sizeLock.Lock()
		defer f.sizeLock.Unlock()
		f.SizeBytes += info.Size()
	}

	return nil
}

func (f *File) directorySize() error {
	return filepath.Walk(f.Path, f.walkFn)
}

// GetField returns the field value for this file.
func (f *File) GetField(field Field) interface{} {
	switch field {
	case Name:
		return f.Name
	case Extension:
		return f.Extension
	case Path:
		return f.Path
	case Age:
		return f.Age
	case Size:
		return f.SizeBytes
	default:
		panic("NOT IMPLEMENTED")
	}
}
