package task

import (
    "github.com/hnakamur/ringbuffer"
)

type Runner interface {
    Run() error
}

type TaskRunner struct {
    queue *ringbuffer.RingBuffer
}

func (r *TaskRunner) Add(tasks ...Runner) error {
    for _, task := range tasks {
        err := r.queue.Add(task)
        if err != nil {
            return err
        }
    }
    return nil
}

func (r *TaskRunner) HasTask() bool {
    return r.queue.Len() > 0
}

func (r *TaskRunner) RunOneTask() error {
    item, err := r.queue.Remove()
    if err != nil {
        return err
    }
    task := item.(Runner)
    return task.Run()
}

func (r *TaskRunner) RunLoop() error {
    for r.HasTask() {
        err := r.RunOneTask()
        if err != nil {
            return err
        }
    }
    return nil
}
