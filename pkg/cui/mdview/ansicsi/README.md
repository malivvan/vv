# ansicsi

[![PkgGoDev](https://pkg.go.dev/badge/github.com/malivvan/vv/pkg/cui/mdview/ansicsi)](https://pkg.go.dev/github.com/malivvan/vv/pkg/cui/mdview/ansicsi)
[![codecov](https://codecov.io/gh/pgavlin/ansicsi/branch/master/graph/badge.svg)](https://codecov.io/gh/pgavlin/ansicsi)
[![Go Report Card](https://goreportcard.com/badge/github.com/malivvan/vv/pkg/cui/mdview/ansicsi)](https://goreportcard.com/report/github.com/malivvan/vv/pkg/cui/mdview/ansicsi)
[![Test](https://github.com/malivvan/vv/pkg/cui/mdview/ansicsi/workflows/Test/badge.svg)](https://github.com/malivvan/vv/pkg/cui/mdview/ansicsi/actions?query=workflow%3ATest)

ansicsi provides a Go package that decodes and encodes ANSI control sequences as defined in ECMA-48/ANSI X3.64.

The high-level decoder currently only supports the Set Graphics Rendition control function. All other control
functions are returned as a tuple of (parameter bytes, intermediate bytes, final byte).

The decoder can be called in a loop in order to separate control sequences from normal text:

```go
for len(bytes) > 0 {
	if cmd, size := Decode(bytes); size > 0 {
		switch cmd := cmd.(type) {
			// Handle control functions here
		}
		bytes = bytes[size:]
		continue
	}

	// Handle plain text here

	bytes = bytes[1:]
}
```

A command can be encoded using its Encode method:

```go
resetCommand := SetGraphicsRendition{Command: SGRReset}
sz, err := resetCommand.Encode(w)
```
