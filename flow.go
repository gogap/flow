package flow

import (
	"fmt"
	"github.com/gogap/config"
	"sync"

	"github.com/gogap/context"
	"github.com/gogap/flow/cache"
)

var (
	defaultFlow = New()
)

type Flow struct {
	cache                 cache.Cache
	handlers              map[string]HandlerFunc
	registerdHandlerNames []string

	lock sync.RWMutex
}

type FlowTrans struct {
	flow    *Flow
	ctx     context.Context
	firstFn HandlerFunc
	err     error
	conf    config.Configuration
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

func (p *Flow) Begin(ctx context.Context) *FlowTrans {
	return &FlowTrans{flow: p, ctx: ctx}
}

func (p *Flow) Run(name string, ctx context.Context) (err error) {

	p.lock.RLock()
	h, exist := p.handlers[name]
	p.lock.RUnlock()

	if !exist {
		err = fmt.Errorf("handler %s not exist", name)
		return
	}

	err = h.Run(ctx)

	return
}

func RegisterHandler(name string, handler HandlerFunc) (err error) {
	return defaultFlow.RegisterHandler(name, handler)
}

func ListHandlers() []string {
	return defaultFlow.ListHandlers()
}

func Run(name string, ctx context.Context) (err error) {
	return defaultFlow.Run(name, ctx)
}

func Begin(ctx context.Context) *FlowTrans {
	return defaultFlow.Begin(ctx)
}

func (p *FlowTrans) WithConfig(name string, opts ...config.Option) *FlowTrans {
	if p.ctx == nil {
		p.ctx = context.NewContext()
	}

	p.ctx.WithValue(ctxConfigKey{name}, config.NewConfig(opts...))

	return p
}

func (p *FlowTrans) WithCache(name string, opts ...cache.Option) *FlowTrans {
	if p.ctx == nil {
		p.ctx = context.NewContext()
	}

	c, err := cache.NewCache(name, opts...)
	if err != nil {
		p.err = err
		return p
	}

	p.ctx.WithValue(ctxCacheKey{name}, c)

	return p
}

func (p *FlowTrans) Then(name string) *FlowTrans {

	if p.err != nil {
		return p
	}

	h, exist := p.flow.handlers[name]

	if !exist {
		p.err = fmt.Errorf("handler %s not exist", name)
		return p
	}

	if p.firstFn == nil {
		p.firstFn = h
		return p
	}

	p.firstFn = p.firstFn.Then(h)

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

	return p.firstFn.Run(ctx)
}
