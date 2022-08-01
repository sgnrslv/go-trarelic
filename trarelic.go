package trarelic

import (
	"net/http"
	"os"

	"github.com/opentracing/opentracing-go"
)

type Trarelic struct {
	Tracer  opentracing.Tracer
	NewSpan bool
	Type    string
	Caller  string
	Postfix string
}

type TrarelicOption func(*Trarelic)

// NewTrarelic is a constructor function for *Trarelic.
func NewTrarelic(opts ...TrarelicOption) *Trarelic {
	caller, err := os.Executable()
	if err != nil {
		caller = os.Args[0]
	}

	t := &Trarelic{
		Tracer:  opentracing.GlobalTracer(),
		NewSpan: true,
		Type:    "background",
		Caller:  caller,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// WithTracer defines certain tracer implementation, if needed,
// defaults to opentracing.GlobalTracer().
func WithTracer(tracer opentracing.Tracer) TrarelicOption {
	return func(t *Trarelic) {
		t.Tracer = tracer
	}
}

// WithNewSpan defines whether new span for external request is needed,
// defaults to true.
func WithNewSpan(new bool) TrarelicOption {
	return func(t *Trarelic) {
		t.NewSpan = new
	}
}

// WithType is a newrelic style segmentation by type (background / web),
// defaults to background.
func WithType(tp string) TrarelicOption {
	return func(t *Trarelic) {
		t.Type = tp
	}
}

// WithCaller defines caller of the external request,
// defaults to path of binary executable.
func WithCaller(caller string) TrarelicOption {
	return func(t *Trarelic) {
		t.Caller = caller
	}
}

// WithPostfix redefines new span name to custom one by adding postfix to external url host.
// Can be useful in two cases:
// 1. we need to segregate external requests within one external host, e.g. payout / payout status, etc.
// 2. external host is not human-readable and custom name is more preferable, e.g. ip
func WithPostfix(postfix string) TrarelicOption {
	return func(t *Trarelic) {
		t.Postfix = postfix
	}
}

// GetSpanFromRequest returns span to witch trarelic's tags will be applied.
func (t *Trarelic) GetSpanFromRequest(req *http.Request) opentracing.Span {
	parentSpan := opentracing.SpanFromContext(req.Context())

	if !t.NewSpan {
		return parentSpan
	}

	operationName := req.URL.Host
	if t.Postfix != "" {
		operationName = operationName + " " + t.Postfix
	}

	var opts []opentracing.StartSpanOption
	if parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	span := t.Tracer.StartSpan(operationName, opts...)

	return span
}
