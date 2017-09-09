package flow

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/gogap/config"
)

type PipeTaskRunner struct {
	name       string
	errHandler ErrorHandler
	conf       config.Configuration

	instanceMap *treemap.Map
}

func init() {
	RegisterRunner("PipeTaskRunner", NewPipeTaskRunner)
}

func NewPipeTaskRunner(conf config.Configuration) (runner TaskRunner, err error) {

	pipRunner := &PipeTaskRunner{
		conf:        conf,
		instanceMap: treemap.NewWithStringComparator(),
	}

	runner = pipRunner

	return
}

func (p *PipeTaskRunner) Name() string {
	return fmt.Sprintf("<PipeTaskRunner.%s>", p.name)
}

func (p *PipeTaskRunner) Run(task *Task) {

	for i := 0; i < len(task.steps); i++ {
		step := &task.steps[i]
		handler, err := p.getHandler(step)
		if err != nil {
			task.appendError(err)
			if p.thowError(err, step, task.ctx) {
				continue
			}
			return
		}

		err = handler.Handle(*step, task.ctx)
		if err != nil {
			task.appendError(err)
			if p.thowError(err, step, task.ctx) {
				continue
			}
			return
		}

	}

	return
}

func (p *PipeTaskRunner) thowError(err error, step *Step, ctx Context) bool {
	if err != nil {
		if p.errHandler != nil {
			return p.errHandler.Handle(err, step, ctx)
		}
		return false
	}

	return true
}

func (p *PipeTaskRunner) SetErrorHandler(handler ErrorHandler) {
	p.errHandler = handler
}

func (p *PipeTaskRunner) getHandler(step *Step) (handler Handler, err error) {

	key := fmt.Sprintf("%s.%s.%s", step.Flow(), step.Name(), step.Handler())

	conf := config.NewConfig()

	conf.WithFallback(p.conf.GetConfig(step.Handler())).
		WithFallback(step.conf)

	singleton := conf.GetBoolean("singleton", true)

	if singleton {
		if h, exist := p.instanceMap.Get(key); exist {
			handler = h.(Handler)
			return
		}
	}

	handler, err = NewHandler(step.Handler(), conf)

	if err == nil {
		p.instanceMap.Put(key, handler)
	}

	return

}
