## GO-trarelic

# Prerequisites
Opentracing tracer must be instantiated and achievable via `opentracing.GlobalTracer()` or passed directly via option `trarelic.WithTracer()`.

# Description
Trarelic creates opentracing span with some tags needed in pulsometer:
- `is_external`
- `type`
- `caller`

You must instrument your code by wrapping the http transport used in external requests.

# E.g.:

Wrap existing transport
```go
client := &http.Client{
	Timeout: 10 * time.Second,
	Transport: trarelic.NewRoundTripper(
		&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	),
}
```

http.DefaultTransport will be used by default if nil is passed
```go
client := &http.Client{
	Transport: trarelic.NewRoundTripper(nil),
}
```

With functional options 
```go
client := &http.Client{
	Transport: trarelic.NewRoundTripper(
		nil,
		trarelic.WithTracer(opentracing.GlobalTracer()),
		trarelic.WithNewSpan(true),
		trarelic.WithType("background"),
		trarelic.WithCaller("bin/i_am_bin"),
		trarelic.WithPostfix("additional info about the call"),
	),
}
```

See [godoc](https://pkg.go.dev/github.com/sgnrslv/go-trarelic) for more information.