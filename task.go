package flow

import (
	"fmt"
	"github.com/pborman/uuid"
	"sync"
)

type TaskStatus int

var (
	TaskStatusReady   TaskStatus = 1
	TaskStatusPending TaskStatus = 2
	TaskStatusRunning TaskStatus = 3
	TaskStatusDone    TaskStatus = 4
)

type Task struct {
	id       string
	flowName string

	steps []Step

	ctx    Context
	runner TaskRunner

	status       TaskStatus
	statusLocker sync.RWMutex

	errs []error

	locked bool
}

func NewTask(flo *Flow, ctx Context) *Task {
	t := &Task{
		id:       uuid.New(),
		runner:   flo.runner,
		flowName: flo.name,
		ctx:      ctx,
		steps:    flo.steps,
	}

	t.status = TaskStatusReady

	return t
}

func (p *Task) Context() Context {
	return p.ctx
}

func (p *Task) LatestError() error {
	if len(p.errs) == 0 {
		return nil
	}

	return p.errs[len(p.errs)-1]
}

func (p *Task) Errors() []error {
	return p.errs
}

func (p *Task) appendError(err error) {
	p.errs = append(p.errs, err)
}

func (p *Task) ID() string {
	return p.id
}

func (p *Task) Flow() string {
	return p.flowName
}

func (p *Task) String() string {
	return fmt.Sprintf("<%s.%s>", p.flowName, p.id)
}

func (p *Task) Status() TaskStatus {
	p.statusLocker.RLock()
	defer p.statusLocker.RUnlock()

	return p.status
}

func (p *Task) Run() (err error) {

	if len(p.errs) > 0 {
		return p.LatestError()
	}

	if p.runner == nil {
		return
	}

	if p.Status() != TaskStatusReady {
		err = fmt.Errorf("flow %s of task %s is not ready, current status is %s", p.flowName, p.id, toStringState(p.status))
		return
	}

	p.statusLocker.Lock()
	defer p.statusLocker.Unlock()

	p.locked = true

	p.status = TaskStatusPending

	p.status = TaskStatusRunning

	p.runner.Run(p)

	p.status = TaskStatusDone

	err = p.LatestError()

	return
}

func toStringState(status TaskStatus) string {
	switch status {
	case TaskStatusReady:
		return "Ready"
	case TaskStatusPending:
		return "Pending"
	case TaskStatusRunning:
		return "Running"
	case TaskStatusDone:
		return "Done"
	default:
		return "Unknown"
	}
}
