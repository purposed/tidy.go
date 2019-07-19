package fsclean

import (
	"github.com/sirupsen/logrus"
)

// Rule maps a condition to an action.
type Rule struct {
	Name      string    `json:"name,omitempty"`
	Condition Condition `json:"condition,omitempty"`

	Action Action `json:"action,omitempty"`

	log *logrus.Entry
}

func (r *Rule) initLogger() {
	if r.log == nil {
		r.log = logrus.WithField("rule", r.Name)
	}
}

// Apply conditionnaly applies the action.
func (r *Rule) Apply(file *File) error {
	r.initLogger()

	ret, err := r.Condition.Evaluate(file)
	if err != nil {
		return err
	}

	if ret {
		r.log.Infof("[%s] - %s", r.Action.Name(), file.Name)
		return r.Action.Execute(file)
	}
	return nil
}
