package runtime

import "fmt"

type Request struct {
	In []byte
}

type Response struct {
	Out []byte
}

type Func struct {
	Name string
	Call func(req *Request) (*Response, error)
}

func (f *Func) String() string {
	return fmt.Sprintf("func: %s\n", f.Name)
}
