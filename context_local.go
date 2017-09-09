package flow

import (
	"github.com/emirpasic/gods/maps/treemap"
)

type LocalContext struct {
	id  string
	ctx *treemap.Map
}

func (p *LocalContext) ID() string {
	return p.id
}

func (p *LocalContext) Get(key string) (value interface{}, exist bool) {
	return p.ctx.Get(key)
}

func (p *LocalContext) Set(key string, value interface{}) Context {
	p.ctx.Put(key, value)
	return p
}

func (p *LocalContext) Delete(key string) Context {
	p.ctx.Remove(key)
	return p
}

func (p *LocalContext) Flush() {
	p.ctx.Clear()
}

func (p LocalContext) Keys() []string {
	var retval []string

	for _, key := range p.ctx.Keys() {
		retval = append(retval, key.(string))
	}

	return retval
}

func (p LocalContext) GetAll() map[string]interface{} {
	var retval = make(map[string]interface{})
	p.ctx.All(func(k, v interface{}) bool {
		retval[k.(string)] = v
		return true
	})

	return retval
}
