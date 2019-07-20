package tidy

import "time"

// File represents a file on disk.
type File struct {
	Name        string
	Path        string
	IsDirectory bool
	Extension   string
	Age         time.Duration
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
	default:
		panic("NOT IMPLEMENTED")
	}
}
