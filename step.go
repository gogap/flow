package flow

import (
	"fmt"
	"github.com/gogap/config"
)

type Step struct {
	name    string
	handler string

	flowName string
	taskId   string

	conf config.Configuration
}

func NewStep(flowName, stepName, handlerName string, conf config.Configuration) Step {
	return Step{
		flowName: flowName,
		name:     stepName,
		handler:  handlerName,
		conf:     conf,
	}
}

func (p *Step) Name() string {
	return p.name
}

func (p *Step) Flow() string {
	return p.flowName
}

func (p *Step) TaskID() string {
	return p.taskId
}

func (p *Step) Handler() string {
	return p.handler
}

func (p *Step) String() string {
	return fmt.Sprintf("<%s.%s.%s.%s>", p.flowName, p.name, p.handler, p.taskId)
}
