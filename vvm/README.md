
# The vv Language

[![GoDoc](https://godoc.org/github.com/malivvan/vv?status.svg)](https://godoc.org/github.com/malivvan/vv)
![test](https://github.com/malivvan/vv/workflows/test/badge.svg)

**vvm is a small, dynamic, fast, secure script language for Go.** 

vvm is **[fast](#benchmark)** and secure because it's compiled/executed as
bytecode on stack-based VM that's written in native Go.

```golang
/* The vvm Language */
fmt := import("fmt")

each := func(seq, fn) {
    for x in seq { fn(x) }
}

sum := func(init, seq) {
    each(seq, func(x) { init += x })
    return init
}

fmt.println(sum(0, [1, 2, 3]))   // "6"
fmt.println(sum("", [1, 2, 3]))  // "123"
```

## Features

- Simple and highly readable
  [Syntax](https://github.com/malivvan/vv/blob/master/docs/tutorial.md)
  - Dynamic typing with type coercion
  - Higher-order functions and closures
  - Immutable values
- [Securely Embeddable](https://github.com/malivvan/vv/blob/master/docs/interoperability.md)
  and [Extensible](https://github.com/malivvan/vv/blob/master/docs/objects.md)
- Compiler/runtime written in native Go _(no external deps or cgo)_
- Executable as a
  [standalone](https://github.com/malivvan/vv/blob/master/docs/vvm-cli.md)
  language / REPL
- Use cases: rules engine, [state machine](https://github.com/d5/go-fsm),
  data pipeline, [transpiler](https://github.com/d5/tengo2lua)

## Benchmark

|                                                      | fib(35) | fibt(35) |  Language (Type)  |
|:-----------------------------------------------------|    ---: |     ---: |  :---: |
| [**vv**](https://github.com/malivvan/vv)                  | `2,931ms` | `4ms` | vvm (VM) |
| [go-lua](https://github.com/Shopify/go-lua)          | `4,824ms` | `4ms` | Lua (VM) |
| [GopherLua](https://github.com/yuin/gopher-lua)      | `5,365ms` | `4ms` | Lua (VM) |
| [goja](https://github.com/dop251/goja)               | `5,533ms` | `5ms` | JavaScript (VM) |
| [starlark-go](https://github.com/google/starlark-go) | `11,495ms` | `5ms` | Starlark (Interpreter) |
| [Yaegi](https://github.com/containous/yaegi)         | `15,645ms` | `12ms` | Yaegi (Interpreter) |
| [gpython](https://github.com/go-python/gpython)      | `16,322ms` | `5ms` | Python (Interpreter) |
| [otto](https://github.com/robertkrimen/otto)         | `73,093ms` | `10ms` | JavaScript (Interpreter) |
| [Anko](https://github.com/mattn/anko)                | `79,809ms` | `8ms` | Anko (Interpreter) |
| -                                                    | - | - | - |
| Go                                                   | `53ms` | `3ms` | Go (Native) |
| Lua                                                  | `1,612ms` | `3ms` | Lua (Native) |
| Python                                               | `2,632ms` | `23ms` | Python 2 (Native) |

## Quick Start

```
go get github.com/malivvan/vv/vvm
```

A simple Go example code that compiles/runs vvm script code with some input/output values:

```golang
package main

import (
	"context"
	"fmt"

	"github.com/malivvan/vv/vvm"
)

func main() {
	// vvm script code
	src := `
each := func(seq, fn) {
    for x in seq { fn(x) }
}

sum := 0
mul := 1
each([a, b, c, d], func(x) {
	sum += x
	mul *= x
})`

	// create a new Script instance
	script := vvm.NewScript([]byte(src))

	// set values
	_ = script.Add("a", 1)
	_ = script.Add("b", 9)
	_ = script.Add("c", 8)
	_ = script.Add("d", 4)

	// run the script
	compiled, err := script.RunContext(context.Background())
	if err != nil {
		panic(err)
	}

	// retrieve values
	sum := compiled.Get("sum")
	mul := compiled.Get("mul")
	fmt.Println(sum, mul) // "22 288"
}
```

## References

- [Language Syntax](https://github.com/malivvan/vv/blob/master/docs/tutorial.md)
- [Object Types](https://github.com/malivvan/vv/blob/master/docs/objects.md)
- [Runtime Types](https://github.com/malivvan/vv/blob/master/docs/runtime-types.md)
  and [Operators](https://github.com/malivvan/vv/blob/master/docs/operators.md)
- [Builtin Functions](https://github.com/malivvan/vv/blob/master/docs/builtins.md)
- [Interoperability](https://github.com/malivvan/vv/blob/master/docs/interoperability.md)
- [vv CLI](https://github.com/malivvan/vv/blob/master/docs/vv-cli.md)
- [Standard Library](https://github.com/malivvan/vv/blob/master/docs/stdlib.md)
- Syntax Highlighters: [VSCode](https://github.com/lissein/vscode-vvm), [Atom](https://github.com/malivvan/vv-atom)
- **Why the name vvm?** It's from [1Q84](https://en.wikipedia.org/wiki/1Q84).

##

:hearts: Like writing Go code? Come work at Skool. [We're hiring!](https://jobs.lever.co/skool)

