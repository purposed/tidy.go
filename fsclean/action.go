package fsclean

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
)

const (
	deleteAction = "delete"
	renameAction = "move"
)

// Action represents a file action.
type Action interface {
	Execute(f *File) error
	Name() string
}

// ActionDelete represents a delete action.
type ActionDelete struct{}

// Name names the action.
func (d *ActionDelete) Name() string {
	return deleteAction
}

// Execute executes the delete action.
func (d *ActionDelete) Execute(f *File) error {
	return os.RemoveAll(f.Path)
}

// ActionRename represents a rename action.
type ActionRename struct {
	NameTemplate string `mapstructure:"name_template"`
	ToDirectory  string `mapstructure:"to_directory"`
}

// Name names the action.
func (d *ActionRename) Name() string {
	return renameAction
}

// Execute executes the action.
func (d *ActionRename) Execute(f *File) error {
	parentDir := path.Dir(f.Path)
	if d.ToDirectory != "" {
		parentDir = d.ToDirectory
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return err
		}
	}

	newName := filepath.Base(f.Path)
	if d.NameTemplate != "" {
		newName = strings.ReplaceAll(d.NameTemplate, "{name}", f.Name)
		newName = strings.ReplaceAll(newName, "{extension}", f.Extension)
	}

	newPath := path.Join(parentDir, newName)

	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return os.Rename(f.Path, newPath)
	}
	return fmt.Errorf("file [%s] already exists", newPath)
}

func getAction(def ActionDefinition) (Action, error) {
	switch def.Type {
	case deleteAction:
		return &ActionDelete{}, nil
	case renameAction:
		var ren ActionRename
		if err := mapstructure.Decode(def.Parameters, &ren); err != nil {
			return nil, err
		}
		return &ren, nil
	default:
		return nil, fmt.Errorf("unknown action: %s", def.Type)
	}
}
