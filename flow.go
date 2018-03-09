package flow

import (
	"fmt"
	"sync"

	"github.com/gogap/config"
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
	flow          *Flow
	ctx           context.Context
	firstFn       HandlerFunc
	err           error
	conf          config.Configuration
	defaultParams Params
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

func (p *Flow) Begin(ctx context.Context, defaultParams ...Params) *FlowTrans {

	var params Params

	if len(defaultParams) > 0 {
		params = defaultParams[0]
	}

	return &FlowTrans{flow: p, ctx: ctx, firstFn: voidTransBegin, defaultParams: params}
}

func (p *Flow) Run(name string, ctx context.Context, params ...Params) (err error) {

	p.lock.RLock()
	h, exist := p.handlers[name]
	p.lock.RUnlock()

	if !exist {
		err = fmt.Errorf("handler %s not exist", name)
		return
	}

	if len(params) > 0 {
		err = h.Run(ctx, params[0])
		return
	}

	err = h.Run(ctx, nil)

	return
}

func RegisterHandler(name string, handler HandlerFunc) (err error) {
	return defaultFlow.RegisterHandler(name, handler)
}

func ListHandlers() []string {
	return defaultFlow.ListHandlers()
}

func Run(name string, ctx context.Context, params ...Params) (err error) {
	return defaultFlow.Run(name, ctx, params...)
}

func Begin(ctx context.Context, defaultParams ...Params) *FlowTrans {
	return defaultFlow.Begin(ctx, defaultParams...)
}

func (p *FlowTrans) Then(name string, params ...Params) *FlowTrans {

	if p.err != nil {
		return p
	}

	h, exist := p.flow.handlers[name]

	if !exist {
		p.err = fmt.Errorf("handler %s not exist", name)
		return p
	}

	var nextParam Params

	if len(params) == 0 {
		nextParam = p.defaultParams
	} else {
		nextParam = params[0]
	}

	p.firstFn = p.firstFn.Then(h, nextParam)

	return p
}

func (p *FlowTrans) Subscribe(subscribers ...SubscriberFunc) *FlowTrans {

	if p.err != nil {
		return p
	}

	if p.firstFn == nil {
		return p
	}

	p.firstFn = p.firstFn.Subscribe(subscribers...)

	return p
}

func (p *FlowTrans) Commit() error {
	if p.err != nil {
		return p.err
	}

	var ctx context.Context
	if p.ctx == nil {
		ctx = context.NewContext()
	} else {
		ctx = p.ctx
	}

	return p.firstFn.Run(ctx, nil)
}

func voidTransBegin(ctx context.Context, params Params) error {
	return nil
}
