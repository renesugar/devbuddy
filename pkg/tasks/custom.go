package tasks

import (
	"fmt"
)

func init() {
	allTasks["custom"] = newCustom
}

type Custom struct {
	taskName  string
	condition string
	command   string
}

func newCustom(config *taskConfig) (Task, error) {
	task := &Custom{}

	properties := config.payload.(map[interface{}]interface{})

	name, ok := properties["name"]
	if !ok {
		name = ""
	}
	command, ok := properties["meet"]
	if !ok {
		return nil, fmt.Errorf("missing key 'meet'")
	}
	condition, ok := properties["met?"]
	if !ok {
		return nil, fmt.Errorf("missing key 'met?'")
	}

	var err error
	task.taskName, err = asString(name)
	if err != nil {
		return nil, fmt.Errorf("invalid name value: %s", err)
	}
	task.command, err = asString(command)
	if err != nil {
		return nil, fmt.Errorf("invalid meet value: %s", err)
	}
	task.condition, err = asString(condition)
	if err != nil {
		return nil, fmt.Errorf("invalid met? value: %s", err)
	}

	return task, nil
}

func (c *Custom) name() string {
	return "Custom"
}

func (c *Custom) header() string {
	if c.taskName != "" {
		return c.taskName
	}
	return c.command
}

func (c *Custom) actions(ctx *context) []taskAction {
	return []taskAction{
		&customAction{condition: c.condition, command: c.command},
	}
}

type customAction struct {
	condition string
	command   string
}

func (c *customAction) description() string {
	return ""
}

func (c *customAction) needed(ctx *context) (bool, error) {
	code, err := shellSilent(ctx, c.condition).RunWithCode()
	if err != nil {
		return false, fmt.Errorf("failed to run the condition command: %s", err)
	}
	return code != 0, nil
}

func (c *customAction) run(ctx *context) error {
	err := shell(ctx, c.command).Run()
	if err != nil {
		return fmt.Errorf("command failed: %s", err)
	}
	return nil
}
