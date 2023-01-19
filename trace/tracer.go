package trace

// Interface that describes an object capable of tracing events throughout code
type Tracer interface {
	Trace(...interface{})
}
