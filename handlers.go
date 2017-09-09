package flow

import (
	"fmt"
	"github.com/gogap/config"
	"sort"
	"sync"
)

type Handler interface {
	Handle(step Step, ctx Context) error
}

type ErrorHandler interface {
	Handle(err error, step *Step, ctx Context) bool
}

type NewHandlerFunc func(conf config.Configuration) (handler Handler, err error)

var (
	handlersLocker  = sync.Mutex{}
	newHandlerFuncs = make(map[string]NewHandlerFunc)
)

func RegisterHandler(name string, newHandlerFunc NewHandlerFunc) {
	handlersLocker.Lock()
	defer handlersLocker.Unlock()

	if name == "" {
		panic("flow: Register handler name is empty")
	}

	if newHandlerFunc == nil {
		panic("flow: Register handler is nil")
	}

	if _, exist := newHandlerFuncs[name]; exist {
		panic("flow: Register called twice for handler " + name)
	}

	newHandlerFuncs[name] = newHandlerFunc
}

func Handlers() []string {
	handlersLocker.Lock()
	defer handlersLocker.Unlock()

	var list []string
	for name := range newHandlerFuncs {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func NewHandler(name string, conf config.Configuration) (handler Handler, err error) {
	handlersLocker.Lock()
	defer handlersLocker.Unlock()

	newHandlerFunc, exist := newHandlerFuncs[name]

	if !exist {
		err = fmt.Errorf("flow: handler %s not registered", name)
		return
	}

	handler, err = newHandlerFunc(conf)

	return
}
