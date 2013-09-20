package service

import (
	"strings"

	"github.com/hnakamur/commango/task"
)

type State string

const (
	Started   = State("started")
	Stopped   = State("stopped")
	Restarted = State("restarted")
	Reloaded  = State("reloaded")
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

	result, err = task.DoRun(func(result *task.Result) error {
		result.Module = "service"
		result.Op = "change_state"
		result.Extra["name"] = s.Name
		result.Extra["state"] = string(state)

		var op string
		switch s.State {
		case Started:
			if oldState == Stopped {
				op = "start"
			}
		case Stopped:
			if oldState == Started {
				op = "stop"
			}
		case Restarted:
			if oldState == Started {
				op = "restart"
			} else {
				op = "start"
			}
		case Reloaded:
			if oldState == Started {
				op = "reload"
			} else {
				op = "start"
			}
		}

		if op == "" {
			result.Skipped = true
			return nil
		} else {
			return result.ExecCommand("service", s.Name, op)
		}
	})
	return
}

func (s *Service) state() (state State, err error) {
	_, err = task.DoRun(func(result *task.Result) error {
		result.Module = "service"
		result.Op = "state"
		result.Extra["name"] = s.Name
		err = result.ExecCommand("service", s.Name, "status")
		result.Changed = false
		if result.Rc == 3 {
			state = Stopped
			result.Err = nil
			err = nil
			result.Failed = false
		} else if result.Rc == 0 {
			state = Started
		}
		return err
	})
	return
}

func (s *Service) ensureAutoStart(enabled bool) (result *task.Result, err error) {
	oldEnabled, err := s.autoStartEnabled()
	if err != nil {
		return
	}

	result, err = task.DoRun(func(result *task.Result) error {
		result.Module = "service"
		result.Op = "auto_start"
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
			return nil
		} else {
			return result.ExecCommand("chkconfig", s.Name, op)
		}
	})
	return
}

func (s *Service) autoStartEnabled() (enabled bool, err error) {
	result, err := task.DoRun(func(result *task.Result) error {
		result.Module = "service"
		result.Op = "check_auto_start"
		result.Extra["name"] = s.Name
		err := result.ExecCommand("chkconfig", s.Name, "--list")
		result.Changed = false
		return err
	})
	enabled = strings.Contains(result.Stdout, "\t2:on\t")
	return
}
