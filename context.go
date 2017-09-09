package flow

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/gogap/config"
	"github.com/pborman/uuid"
	"sort"
	"sync"
)

type Context interface {
	ID() string

	Get(key string) (value interface{}, exist bool)
	Set(key string, value interface{}) Context
	Delete(key string) Context
	Flush()

	Keys() []string

	GetAll() map[string]interface{}
}

type ContextProvider interface {
	NewContext(conf config.Configuration) Context
}

type LocalContextProvider struct {
}

func init() {
	RegisterContextProvider("LocalContextProvider", NewLocalContextProvider)
}

func NewLocalContextProvider(conf config.Configuration) (provider ContextProvider, err error) {
	return &LocalContextProvider{}, nil
}

func (p *LocalContextProvider) NewContext(conf config.Configuration) Context {
	locCtx := &LocalContext{
		id:  uuid.New(),
		ctx: treemap.NewWithStringComparator(),
	}

	return locCtx
}

type NewContextProviderFunc func(conf config.Configuration) (contextProvider ContextProvider, err error)

var (
	contextProvidersLocker  = sync.Mutex{}
	newContextProviderFuncs = make(map[string]NewContextProviderFunc)
)

func RegisterContextProvider(name string, newContextProviderFunc NewContextProviderFunc) {
	contextProvidersLocker.Lock()
	defer contextProvidersLocker.Unlock()

	if name == "" {
		panic("flow: Register context provider name is empty")
	}

	if newContextProviderFunc == nil {
		panic("flow: Register context provider is nil")
	}

	if _, exist := newContextProviderFuncs[name]; exist {
		panic("flow: Register called twice for context provider " + name)
	}

	newContextProviderFuncs[name] = newContextProviderFunc
}

func ContextProviders() []string {
	contextProvidersLocker.Lock()
	defer contextProvidersLocker.Unlock()

	var list []string
	for name := range newContextProviderFuncs {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func NewContextProvider(name string, conf config.Configuration) (contextProvider ContextProvider, err error) {
	contextProvidersLocker.Lock()
	defer contextProvidersLocker.Unlock()

	newContextProviderFunc, exist := newContextProviderFuncs[name]

	if !exist {
		err = fmt.Errorf("flow: context provider driver %s not registered", name)
		return
	}

	contextProvider, err = newContextProviderFunc(conf)

	return
}
