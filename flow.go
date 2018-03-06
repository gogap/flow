package flow

import (
	"fmt"
	"sync"

	"github.com/gogap/context"
)

var (
	defaultFlow = New()
)

type Flow struct {
	handlers              map[string]HandlerFunc
	registerdHandlerNames []string

	lock sync.RWMutex
}

type FlowTrans struct {
	flow *Flow
	h    HandlerFunc
	opts []Option
	err  error
}

func New() *Flow {
	return &Flow{
		handlers: make(map[string]HandlerFunc),
	}
}

func (p *Flow) RegisterHandler(name string, handler HandlerFunc) (err error) {

	if len(name) == 0 {
		err = fmt.Errorf("handler name is emtpy")
		return
	}

	if handler == nil {
		err = fmt.Errorf("handler func could not be nil")
		return
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	_, exist := p.handlers[name]
	if exist {
		err = fmt.Errorf("handler %s already registered", name)
		return
	}

	p.handlers[name] = handler
	p.registerdHandlerNames = append(p.registerdHandlerNames, name)

	return
}

func (p *Flow) ListHandlers() []string {
	return p.registerdHandlerNames
}

func (p *Flow) Begin() *FlowTrans {
	return &FlowTrans{flow: p}
}

func (p *Flow) Run(name string, ctx context.Context, opts ...Option) (err error) {

	p.lock.RLock()
	h, exist := p.handlers[name]
	p.lock.RUnlock()

	if !exist {
		err = fmt.Errorf("handler %s not exist", name)
		return
	}

	err = h.Run(ctx, opts...)

	return
}

func RegisterHandler(name string, handler HandlerFunc) (err error) {
	return defaultFlow.RegisterHandler(name, handler)
}

func ListHandlers() []string {
	return defaultFlow.ListHandlers()
}

func Begin() *FlowTrans {
	return &FlowTrans{flow: defaultFlow}
}

func Run(name string, ctx context.Context, opts ...Option) (err error) {
	return defaultFlow.Run(name, ctx, opts...)
}

func (p *FlowTrans) Then(name string, opts ...Option) *FlowTrans {

	if p.err != nil {
		return p
	}

	h, exist := p.flow.handlers[name]

	if !exist {
		p.err = fmt.Errorf("handler %s not exist", name)
		return p
	}

	if p.h == nil {
		p.h = h
		p.opts = opts
		return p
	}

	p.h = p.h.Then(h, opts...)

	return p
}

func (p *FlowTrans) Subscribe(subscribers ...SubscriberFunc) *FlowTrans {

	if p.err != nil {
		return p
	}

	if p.h == nil {
		return p
	}

	p.h = p.h.Subscribe(subscribers...)

	return p
}

func (p *FlowTrans) Commit() error {
	if p.err != nil {
		return p.err
	}

	return p.h.Run(context.NewContext(), p.opts...)
}
