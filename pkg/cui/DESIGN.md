This document lists architectural details of cui.

# Focus-related style attributes are unset by default

This applies to all widgets except Button and TabbedPanels, which require a
style change to indicate focus. See [ColorUnset](https://docs.rocket9labs.com/github.com/malivvan/vv/pkg/cui#pkg-variables).

# Widgets always use `sync.RWMutex`

See [#30](https://github.com/malivvan/vv/pkg/cui/issues/30).
