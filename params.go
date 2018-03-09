package flow

import (
	"fmt"
	"time"
)

type Params map[interface{}]interface{}

func (p Params) IsEmpty() bool {
	return p == nil || len(p) == 0
}

func (p Params) Val(key interface{}) interface{} {
	if p.IsEmpty() {
		return nil
	}

	v := p[key]

	return v
}

func (p Params) Exist(key interface{}) bool {
	if p.IsEmpty() {
		return false
	}

	_, exist := p[key]

	return exist
}

func (p Params) String(key interface{}) string {
	if p == nil {
		return ""
	}

	v, exist := p[key]
	if !exist {
		return ""
	}

	val, ok := v.(string)

	if ok {
		return val
	}

	return fmt.Sprintf("%v", val)
}

func (p Params) Int(key interface{}) int {

	if p == nil {
		return 0
	}

	v, exist := p[key]
	if !exist {
		return 0
	}

	val, ok := v.(int)

	if ok {
		return val
	}

	return 0
}

func (p Params) Int32(key interface{}) int32 {

	if p == nil {
		return 0
	}

	v, exist := p[key]
	if !exist {
		return 0
	}

	val, ok := v.(int32)

	if ok {
		return val
	}

	return 0
}

func (p Params) Int64(key interface{}) int64 {

	if p == nil {
		return 0
	}

	v, exist := p[key]
	if !exist {
		return 0
	}

	val, ok := v.(int64)

	if ok {
		return val
	}

	return 0
}

func (p Params) Float32(key interface{}) float32 {

	if p == nil {
		return 0
	}

	v, exist := p[key]
	if !exist {
		return 0
	}

	val, ok := v.(float32)

	if ok {
		return val
	}

	return 0
}

func (p Params) Float64(key interface{}) float64 {

	if p == nil {
		return 0
	}

	v, exist := p[key]
	if !exist {
		return 0
	}

	val, ok := v.(float64)

	if ok {
		return val
	}

	return 0
}

func (p Params) Duration(key interface{}) time.Duration {

	if p == nil {
		return 0
	}

	v, exist := p[key]
	if !exist {
		return 0
	}

	val, ok := v.(time.Duration)

	if ok {
		return val
	}

	return 0
}

func (p Params) Boolean(key interface{}) bool {

	if p == nil {
		return false
	}

	v, exist := p[key]
	if !exist {
		return false
	}

	val, ok := v.(bool)

	if ok {
		return val
	}

	return false
}
