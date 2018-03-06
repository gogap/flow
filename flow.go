package flow

import (
	"fmt"
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
	flow        *Flow
	ctx         context.Context
	firstFn     HandlerFunc
	firstOpts   []Option
	defaultOpts []Option
	err         error
}

func New() *Flow {
	return &Flow{
		handlers: make(map[string]HandlerFunc),
	}
}

func (p *Flow) WithCache(cache cache.Cache) *Flow {
	p.cache = cache
	return p
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

func WithCache(cache cache.Cache) *Flow {
	return defaultFlow.WithCache(cache)
}

func Run(name string, ctx context.Context, opts ...Option) (err error) {
	return defaultFlow.Run(name, ctx, opts...)
}

func Begin() *FlowTrans {
	return &FlowTrans{flow: defaultFlow}
}

func (p *FlowTrans) WithContext(ctx context.Context) *FlowTrans {
	if p.ctx != nil {
		panic("WithContext only could be call once.")
	}
	p.ctx = ctx
	return p
}

func (p *FlowTrans) WithOptions(opts ...Option) *FlowTrans {

	runOpts := ParseOptions(opts...)
	if runOpts.Cache == nil && p.flow.cache != nil {
		opts = append(opts, Cache(p.flow.cache))
	}

	p.defaultOpts = opts
	return p
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

	if len(opts) == 0 {
		opts = p.defaultOpts
	}

	if p.firstFn == nil {
		p.firstFn = h
		p.firstOpts = opts
		return p
	}

	p.firstFn = p.firstFn.Then(h, opts...)

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

	return p.firstFn.Run(ctx, p.firstOpts...)
}
