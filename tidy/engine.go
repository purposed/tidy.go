package tidy

import (
	"time"

	"github.com/purposed/good/task"
	"github.com/sirupsen/logrus"
)

// An Engine enforces all the monitors.
type Engine struct {
	cfg      *Config
	log      *logrus.Entry
	monitors []*Monitor
	tasks    []*task.RecurringTask
}

// NewEngine initializes the engine.
func NewEngine(cfg *Config) (*Engine, error) {
	monitors, err := cfg.GetMonitors()
	if err != nil {
		return nil, err
	}

	return &Engine{
		cfg:      cfg,
		log:      logrus.WithField("component", "engine"),
		monitors: monitors,
	}, nil
}

// Start starts the engine.
func (e *Engine) Start() {
	log := logrus.New()
	i := 0
	for _, mon := range e.monitors {
		newTask := task.New(task.Parameters{Name: mon.RootDirectory, Function: mon.Check, Logger: log})
		e.tasks = append(e.tasks, newTask)
		newTask.Start(time.Duration(mon.CheckFrequencySeconds) * time.Second)
		i++
	}
	e.log.Infof("started the watcher engine with %d monitors", i)
}

// Stop stops the engine.
func (e *Engine) Stop() {
	for _, t := range e.tasks {
		t.Stop()
	}
	e.log.Info("terminated")
}
