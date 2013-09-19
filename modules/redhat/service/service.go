package service

import (
	"strings"

	"github.com/hnakamur/commango/task"
)

type State string

const (
	STARTED   = State("started")
	STOPPED   = State("stopped")
	RESTARTED = State("restarted")
	RELOADED  = State("reloaded")
)

type Service struct {
	State            State
	Name             string
	AutoStartEnabled bool
}

func (s *Service) Run() (*task.Result, error) {
	result, err := s.ensureState(s.State)
	if err != nil {
		return result, err
	}

	return s.ensureAutoStart(s.AutoStartEnabled)
}

func (s *Service) ensureState(state State) (result *task.Result, err error) {
	oldState, err := s.state()
	if err != nil {
		return
	}

	result = task.NewResult("service.change_state")
	result.RecordStartTime()

	result.Extra["name"] = s.Name
	result.Extra["state"] = string(state)

	var op string
	switch s.State {
	case STARTED:
		if oldState == STOPPED {
			op = "start"
		}
	case STOPPED:
		if oldState == STARTED {
			op = "stop"
		}
	case RESTARTED:
		if oldState == STARTED {
			op = "restart"
		} else {
			op = "start"
		}
	case RELOADED:
		if oldState == STARTED {
			op = "reload"
		} else {
			op = "start"
		}
	}

	if op == "" {
		result.Skipped = true
	} else {
        err = result.ExecCommand("service", s.Name, op)
    }
    result.RecordEndTime()
    result.Log()
	return
}

func (s *Service) state() (state State, err error) {
	result := task.NewResult("service.state")
	result.RecordStartTime()

	err = result.ExecCommand("service", s.Name, "status")
	result.Extra["name"] = s.Name
	if result.Rc == 3 {
		state = STOPPED
		result.Err = nil
		err = nil
		result.Failed = false
	} else if result.Rc == 0 {
		state = STARTED
	}

	result.Changed = false
	result.RecordEndTime()
	result.Log()
	return
}

func (s *Service) ensureAutoStart(enabled bool) (result *task.Result, err error) {
	oldEnabled, err := s.autoStartEnabled()
	if err != nil {
		return
	}

	result = task.NewResult("service.change_auto_start")
	result.RecordStartTime()

	result.Extra["name"] = s.Name

	var op string
	if enabled {
		if !oldEnabled {
			op = "on"
		}
	} else {
		if oldEnabled {
			op = "off"
		}
	}
	if op == "" {
		result.Skipped = true
	} else {
        err = result.ExecCommand("chkconfig", s.Name, op)
    }
    result.RecordEndTime()
    result.Log()
	return
}

func (s *Service) autoStartEnabled() (enabled bool, err error) {
	result := task.NewResult("service.auto_start")
	result.RecordStartTime()

	err = result.ExecCommand("chkconfig", s.Name, "--list")
	enabled = strings.Contains(result.Stdout, "\t2:on\t")

	result.Changed = false
	result.RecordEndTime()
	result.Log()
	return
}
