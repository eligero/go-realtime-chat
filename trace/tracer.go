package trace

import (
	"fmt"
	"io"
)

// Interface that describes an object capable of tracing events throughout code
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer // Place to write the trace output
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
