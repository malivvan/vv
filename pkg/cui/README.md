# cview - Terminal-based user interface toolkit
[![GoDoc](https://codeberg.org/tslocum/godoc-static/raw/branch/master/badge.svg)](https://docs.rocket9labs.com/github.com/malivvan/vv/pkg/cui)
[![Donate](https://img.shields.io/liberapay/receives/rocket9labs.com.svg?logo=liberapay)](https://liberapay.com/rocket9labs.com)

This package is a fork of [tview](https://github.com/rivo/tview).
See [FORK.md](https://github.com/malivvan/vv/pkg/cui/src/branch/master/FORK.md) for more information.

## Demo

`ssh cui.rocket9labs.com -p 20000`

[![Recording of presentation demo](https://github.com/malivvan/vv/pkg/cui/raw/branch/master/cui.svg)](https://github.com/malivvan/vv/pkg/cui/src/branch/master/demos/presentation)

## Features

Available widgets:

- __Input forms__ (including __input/password fields__, __drop-down selections__, __checkboxes__, and __buttons__)
- Navigable multi-color __text views__
- Selectable __lists__ with __context menus__
- Modal __dialogs__
- Horizontal and vertical __progress bars__
- __Grid__, __Flexbox__ and __tabbed panel layouts__
- Sophisticated navigable __table views__
- Flexible __tree views__
- Draggable and resizable __windows__
- An __application__ wrapper

Widgets may be customized and extended to suit any application.

[Mouse support](https://docs.rocket9labs.com/github.com/malivvan/vv/pkg/cui#hdr-Mouse_Support) is available.

## Applications

A list of applications powered by cview is available via [pkg.go.dev](https://pkg.go.dev/github.com/malivvan/vv/pkg/cui?tab=importedby).

## Installation

```bash
go get github.com/malivvan/vv/pkg/cui
```

## Hello World

This basic example creates a TextView titled "Hello, World!" and displays it in your terminal:

```go
package main

import (
	"github.com/malivvan/vv/pkg/cui"
)

func main() {
	app := cui.NewApplication()

	tv := cui.NewTextView()
	tv.SetBorder(true)
	tv.SetTitle("Hello, world!")
	tv.SetText("Lorem ipsum dolor sit amet")
	
	app.SetRoot(tv, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
```

Examples are available via [godoc](https://docs.rocket9labs.com/github.com/malivvan/vv/pkg/cui#pkg-examples)
and in the [demos](https://github.com/malivvan/vv/pkg/cui/src/branch/master/demos) directory.

For a presentation highlighting the features of this package, compile and run
the program in the [demos/presentation](https://github.com/malivvan/vv/pkg/cui/src/branch/master/demos/presentation) directory.

## Documentation

Package documentation is available via [godoc](https://docs.rocket9labs.com/github.com/malivvan/vv/pkg/cui).

An [introduction tutorial](https://rocket9labs.com/post/tview-and-you/) is also available.

## Dependencies

This package is based on [github.com/gdamore/tcell](https://github.com/gdamore/tcell)
(and its dependencies) and [github.com/rivo/uniseg](https://github.com/rivo/uniseg).

## Support

[CONTRIBUTING.md](https://github.com/malivvan/vv/pkg/cui/src/branch/master/CONTRIBUTING.md) describes how to share
issues, suggestions and patches (pull requests).

## Packages
- / [codeberg.org/tslocum/cview](https://codeberg.org/tslocum/cview/src/commit/242e7c1f1b61a4b3722a1afb45ca1165aefa9a59)
- /bind.go [codeberg.org/tslocum/cbind](https://codeberg.org/tslocum/cbind/src/commit/5cd49d3cfccbe4eefaab8a5282826aa95100aa42)
- /vte/ [git.sr.ht/~rockorager/tcell-term](https://git.sr.ht/~rockorager/tcell-term/refs/v0.10.0)
- /femto/ [github.com/wellcomez/femto](https://github.com/wellcomez/femto/tree/8413a0288bcb042fd0de5cbbcb9893c16a01ee69)
- /chart/ [github.com/navidys/tvxwidgets](https://github.com/navidys/tvxwidgets/commit/96bcc0450684693eebd4f8e3e95fcc40eae2dbaa)
