package tidy

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

// Monitor monitors a directory and applies rules.
type Monitor struct {
	RootDirectory         string
	Rules                 []*Rule
	Recursive             bool
	CheckFrequencySeconds int

	log *logrus.Entry
}

// NewMonitor initializes & returns a monitor object.
func NewMonitor(rootDir string, rules []*Rule, recursive bool, checkFrequencySeconds int) (*Monitor, error) {
	if rootDir == "" {
		return nil, errors.New("invalid root directory")
	}

	if checkFrequencySeconds == 0 {
		return nil, errors.New("check frequency cannot be 0")
	}

	dir, err := homedir.Expand(rootDir)
	if err != nil {
		return nil, err
	}

	return &Monitor{
		RootDirectory:         dir,
		Rules:                 rules,
		Recursive:             recursive,
		CheckFrequencySeconds: checkFrequencySeconds,
		log:                   logrus.WithField("component", "monitor"),
	}, nil
}

func (m *Monitor) apply(path string, info os.FileInfo) error {
	f, err := NewFile(path, info)
	if err != nil {
		return err
	}

	for _, rule := range m.Rules {
		if err := rule.Apply(&f); err != nil {
			return err
		}
	}
	return nil
}

func (m *Monitor) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	return m.apply(path, info)
}

// Check recursively checks all file in the directory for the rules.
func (m *Monitor) Check() error {
	m.log.Infof("checking monitor for [%s]", m.RootDirectory)
	if m.Recursive {
		return filepath.Walk(m.RootDirectory, m.walkFunc)
	}

	files, err := ioutil.ReadDir(m.RootDirectory)
	if err != nil {
		return err
	}
	for _, info := range files {
		path := path.Join(m.RootDirectory, info.Name())
		if err := m.apply(path, info); err != nil {
			return err
		}
	}
	return nil
}
