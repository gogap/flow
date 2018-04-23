package flow

import (
	"encoding/json"
	"sync"

	"github.com/gogap/context"
)

type outputKey struct{}

type NameValue struct {
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
	Tags  []string        `json:"tags,omitempty"`
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

func (p *Output) Append(items ...NameValue) {

	p.locker.Lock()
	defer p.locker.Unlock()

	output := p
	for output != nil {
		if output.next != nil {
			output = output.next
			continue
		}

		for _, item := range items {
			output.next = &Output{item: item}
			output = output.next
		}

		return
	}
}

func AppendOutput(ctx context.Context, items ...NameValue) {

	if ctx == nil {
		return
	}

	output, ok := ctx.Value(outputKey{}).(*Output)

	if !ok || output == nil {
		output = &Output{}
		output.Append(items...)
		ctx.WithValue(outputKey{}, output.next)
		return
	}

	output.Append(items...)
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

func FindOutput(ctx context.Context, name string, tags ...string) []NameValue {
	outputList := ListOutput(ctx)

	var matched []NameValue

	for i := 0; i < len(outputList); i++ {
		if name != outputList[i].Name {
			continue
		}

		if !isTagsMatched(tags, outputList[i].Tags) {
			continue
		}

		matched = append(matched, outputList[i])
	}

	return matched
}

func isTagsMatched(tagsA []string, tagsB []string) bool {
	mapTagsB := map[string]bool{}
	for i := 0; i < len(tagsB); i++ {
		mapTagsB[tagsB[i]] = true
	}

	matched := 0

	for i := 0; i < len(tagsA); i++ {
		if mapTagsB[tagsA[i]] == true {
			matched++
		}
	}

	return matched == len(tagsA)
}
