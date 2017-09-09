package flow

import (
	"github.com/gogap/config"
)

type FlowOption func(*Flow) error

type Flow struct {
	name string

	steps []Step

	runner      TaskRunner
	ctxProvider ContextProvider

	conf         config.Configuration
	configFile   string
	configString string
}

func ConfigFile(filename string) FlowOption {
	return func(f *Flow) error {
		conf := config.NewConfig(config.ConfigFile(filename))
		f.conf.WithFallback(conf)
		return nil
	}
}

func ConfigString(str string) FlowOption {
	return func(f *Flow) error {
		conf := config.NewConfig(config.ConfigString(str))
		f.conf.WithFallback(conf)
		return nil
	}
}

func NewFlow(name string, opts ...FlowOption) (f *Flow, err error) {
	flo := &Flow{
		name: name,
		conf: config.NewConfig(),
	}

	for i := 0; i < len(opts); i++ {
		err = opts[i](flo)
		if err != nil {
			return
		}
	}

	ctxConf := flo.conf.GetConfig("context")
	runnerConf := flo.conf.GetConfig("runner")

	ctxProvider := ctxConf.GetString("provider", "LocalContextProvider")
	runnerType := runnerConf.GetString("type", "PipeTaskRunner")

	provider, err := NewContextProvider(ctxProvider, ctxConf.GetConfig("options"))
	if err != nil {
		return
	}

	runner, err := NewRunner(runnerType, runnerConf.GetConfig("options"))
	if err != nil {
		return
	}

	stepsConf := flo.conf.GetConfig("steps")

	orderList := stepsConf.GetStringList("order")

	for i := 0; i < len(orderList); i++ {
		stepConf := stepsConf.GetConfig(orderList[i])
		step := NewStep(
			name,
			orderList[i],
			stepConf.GetString("handler"),
			stepConf.GetConfig("options"),
		)

		flo.steps = append(flo.steps, step)
	}

	flo.runner = runner
	flo.ctxProvider = provider

	f = flo

	return
}

func (p *Flow) Setup(steps []Step) *Flow {

	p.steps = append(p.steps, steps...)

	return p
}

func (p *Flow) Name() string {
	return p.name
}

func (p *Flow) NewTask() *Task {
	return NewTask(p, p.ctxProvider.NewContext(p.conf))
}
