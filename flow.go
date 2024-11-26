package flow

import (
	"errors"
	"fmt"
)

type DefinedInput struct {
	Name  string
	Value interface{}
	Type  string
	Meta  string
}

type ProcessContext struct {
	Vals map[interface{}]interface{}
}

type Process struct {
	Handler      ProcessHandler
	Name         string
	DefinedInput []DefinedInput // name, value, meta   || parse arguments based on the defined inputs
}

type Flow struct {
	Name      string
	Processes []Process
}

func (f *Flow) Process(c *ProcessContext) error {
	for i := 0; i < len(f.Processes); i++ {
		err := f.Processes[i].Handler(c, f.Processes[i].DefinedInput)
		if err != nil {
			eInfo := fmt.Errorf("index: %d\nprocess: %s", i, f.Processes[i].Name)
			return errors.Join(ErrProcessFailure, eInfo, err)
		}
	}
	return nil
}
