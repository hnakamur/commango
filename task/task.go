package task

import (
    "encoding/json"

    "github.com/hnakamur/ringbuffer"
)

type Task interface {
    Run() error
    json.Marshaler
}

type TaskQueue struct {
    queue *ringbuffer.RingBuffer
}

const QUEUE_SIZE = 64

func NewTaskQueue() *TaskQueue {
    return &TaskQueue{
        queue: ringbuffer.NewRingBuffer(QUEUE_SIZE),
    }
}

func (r *TaskQueue) Add(tasks ...Task) error {
    for _, task := range tasks {
        err := r.queue.Add(task)
        if err != nil {
            return err
        }
    }
    return nil
}

func (r *TaskQueue) HasTask() bool {
    return r.queue.Len() > 0
}

func (r *TaskQueue) RunOneTask() error {
    item, err := r.queue.Remove()
    if err != nil {
        return err
    }
    task := item.(Task)
    return task.Run()
}

func (r *TaskQueue) RunLoop() error {
    for r.HasTask() {
        err := r.RunOneTask()
        if err != nil {
            return err
        }
    }
    return nil
}
