package tidy

// Config aggregates tidy configuration.
type Config struct {
	Monitors []MonitorDefinition `json:"monitors"`
}

// GetMonitors converts the config to usable monitors.
func (c *Config) GetMonitors() ([]*Monitor, error) {
	var monitors []*Monitor
	for _, monDef := range c.Monitors {
		mon, err := monDef.ToMonitor()
		if err != nil {
			return nil, err
		}
		monitors = append(monitors, mon)
	}

	return monitors, nil
}

// MonitorDefinition defines a monitor.
type MonitorDefinition struct {
	RootDirectory         string           `json:"root_directory"`
	Rules                 []RuleDefinition `json:"rules"`
	Recursive             bool             `json:"recursive"`
	CheckFrequencySeconds int              `json:"check_frequency_s"`
}

// ToMonitor converts a monitor definition to a monitor object.
func (m *MonitorDefinition) ToMonitor() (*Monitor, error) {
	var rules []*Rule

	for _, ruleDef := range m.Rules {
		rObj, err := ruleDef.ToRule()
		if err != nil {
			return nil, err
		}

		rules = append(rules, rObj)
	}

	return NewMonitor(m.RootDirectory, rules, m.Recursive, m.CheckFrequencySeconds)
}

// RuleDefinition defines a rule.
type RuleDefinition struct {
	Name      string           `json:"name"`
	Condition string           `json:"condition"`
	Action    ActionDefinition `json:"action"`
}

// ToRule converts the rule def to a rule object.
func (r *RuleDefinition) ToRule() (*Rule, error) {
	cond, err := ParseCondition(r.Condition)
	if err != nil {
		return nil, err
	}

	action, err := r.Action.ToAction()
	if err != nil {
		return nil, err
	}

	return &Rule{
		Name:      r.Name,
		Condition: cond,
		Action:    action,
	}, nil
}

// ActionDefinition defines an action.
type ActionDefinition struct {
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
}

// ToAction converts the definition to the object.
func (d *ActionDefinition) ToAction() (Action, error) {
	return getAction(*d)
}
