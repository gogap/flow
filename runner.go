package flow

import (
	"fmt"
	"github.com/gogap/config"
	"sort"
	"sync"
)

type TaskRunner interface {
	Run(task *Task)
	SetErrorHandler(handler ErrorHandler)
}

type NewRunnerFunc func(conf config.Configuration) (runner TaskRunner, err error)

var (
	runnersLocker  = sync.Mutex{}
	newRunnerFuncs = make(map[string]NewRunnerFunc)
)

func RegisterRunner(name string, newRunnerFunc NewRunnerFunc) {
	runnersLocker.Lock()
	defer runnersLocker.Unlock()

	if name == "" {
		panic("flow: Register runner name is empty")
	}

	if newRunnerFunc == nil {
		panic("flow: Register runner is nil")
	}

	if _, exist := newRunnerFuncs[name]; exist {
		panic("flow: Register called twice for runner " + name)
	}

	newRunnerFuncs[name] = newRunnerFunc
}

func Runners() []string {
	runnersLocker.Lock()
	defer runnersLocker.Unlock()

	var list []string
	for name := range newRunnerFuncs {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func NewRunner(name string, conf config.Configuration) (runner TaskRunner, err error) {
	runnersLocker.Lock()
	defer runnersLocker.Unlock()

	newRunnerFunc, exist := newRunnerFuncs[name]

	if !exist {
		err = fmt.Errorf("flow: runner driver %s not registered", name)
		return
	}

	runner, err = newRunnerFunc(conf)

	return
}
