package flow

import (
	"github.com/gogap/context"
	"sync"
)

type outputKey struct{}

type NameValue struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Tags  []string    `json:"tags"`
}

type Output struct {
	item NameValue
	next *Output

	locker sync.Mutex
}

func (p *Output) List() []NameValue {

	if p == nil {
		return nil
	}

	var nv []NameValue

	output := p
	for output != nil {
		nv = append(nv, NameValue{output.item.Name, output.item.Value, output.item.Tags})
		if output.next != nil {
			output = output.next
			continue
		}
		return nv
	}

	return nil
}

func (p *Output) Append(name string, value interface{}, tags []string) {

	p.locker.Lock()
	defer p.locker.Unlock()

	output := p
	for output != nil {
		if output.next != nil {
			output = output.next
			continue
		}

		output.next = &Output{item: NameValue{name, value, tags}}
		return
	}
}

func AppendOutput(ctx context.Context, name string, value interface{}, tags ...string) {

	if ctx == nil {
		return
	}

	output, ok := ctx.Value(outputKey{}).(*Output)

	if !ok {
		ctx.WithValue(outputKey{}, &Output{item: NameValue{name, value, tags}})
		return
	}

	if output == nil {
		output = &Output{item: NameValue{name, value, tags}}
		ctx.WithValue(outputKey{}, output)
		return
	}

	output.Append(name, value, tags)
}

func ListOutput(ctx context.Context) []NameValue {
	if ctx == nil {
		return nil
	}

	output, ok := ctx.Value(outputKey{}).(*Output)

	if !ok {
		return nil
	}

	return output.List()
}
