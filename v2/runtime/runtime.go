package runtime

import (
	"errors"
	"fmt"

	"github.com/derekparker/trie"
)

type Runtime struct {
	registry map[*trie.Node]*Func
}

func NewRuntime() *Runtime {
	return &Runtime{
		make(map[*trie.Node]*Func),
	}
}

func (r *Runtime) RegisterFunc(handle *trie.Node, f *Func) {
	r.registry[handle] = f

	fmt.Printf("Registered func %v\n!", f)
}

func (r *Runtime) Invoke(req *Request, handle *trie.Node) (*Response, error) {
	v, ok := r.registry[handle]
	if !ok {
		return nil, errors.New("could not find handler")
	}

	fmt.Printf("Invoking func: %s\n", v)

	return v.Call(req)
}

// func (r *Runtime) bg(handle *trie.Node, f *Func) {

// }
