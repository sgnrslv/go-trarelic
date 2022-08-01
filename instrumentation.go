package trarelic

import (
	"net/http"
)

// NewRoundTripper creates http.RoundTripper to instrument external requests.
// The RoundTripper returned creates an external span with tags (or use existed one) before delegating to the original
// RoundTripper provided (or http.DefaultTransport if none is provided).
func NewRoundTripper(original http.RoundTripper, opts ...TrarelicOption) http.RoundTripper {
	settings := NewTrarelic(opts...)

	return roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		// The specification of http.RoundTripper requires that the request is never modified.
		// And though we don't need the modification right now, in the future this could be useful.
		request = cloneRequest(request)

		if nil == original {
			original = http.DefaultTransport
		}

		span := settings.GetSpanFromRequest(request)
		if span != nil {
			span.SetTag("is_external", true)
			span.SetTag("type", settings.Type)
			span.SetTag("caller", settings.Caller)
		}

		response, err := original.RoundTrip(request)

		if span != nil && settings.NewSpan {
			span.Finish()
		}

		return response, err
	})
}

// cloneRequest mimics implementation of
// https://godoc.org/github.com/google/go-github/github#BasicAuthTransport.RoundTrip
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
