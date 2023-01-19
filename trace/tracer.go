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

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off creates a Tracer that will ignore calls to Trace
func Off() Tracer {
	return &nilTracer{}
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
