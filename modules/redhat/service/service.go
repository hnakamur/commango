package service

import (
	"os/exec"
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
        return
    }

    cmd := exec.Command("service", s.Name, op)
    result, err = task.ExecCommand("service.change_state", cmd)
	result.Extra["name"] = s.Name
	result.Extra["state"] = string(state)
	return
}

func (s *Service) state() (state State, err error) {
	cmd := exec.Command("service", s.Name, "status")
    result, err := task.ExecCommand("service.state", cmd)
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
	result.Log()
	return state, err
}

func (s *Service) ensureAutoStart(enabled bool) (result *task.Result, err error) {
    oldEnabled, err := s.autoStartEnabled()
    if err != nil {
        return
    }

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
        return
    }
	cmd := exec.Command("chkconfig", s.Name, op)
    result, err = task.ExecCommand("service.change_auto_start", cmd)
	return
}

func (s *Service) autoStartEnabled() (enabled bool, err error) {
	cmd := exec.Command("chkconfig", s.Name, "--list")
    result, err := task.ExecCommand("service.auto_start", cmd)
	enabled = strings.Contains(result.Stdout, "\t2:on\t")
	result.Changed = false
	result.Log()
	return
}
