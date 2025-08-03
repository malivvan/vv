package cli

import (
	"flag"
	"fmt"
	"math"
	"strconv"
)

// TakesValue returns true of the flag takes a value, otherwise false
func (f *UintFlag) TakesValue() bool {
	return true
}

// GetUsage returns the usage string for the flag
func (f *UintFlag) GetUsage() string {
	return f.Usage
}

// GetCategory returns the category for the flag
func (f *UintFlag) GetCategory() string {
	return f.Category
}

// Apply populates the flag given the flag set and environment
func (f *UintFlag) Apply(set *flag.FlagSet) error {
	// set default value so that environment wont be able to overwrite it
	f.defaultValue = f.Value
	f.defaultValueSet = true

	if val, source, found := flagFromEnvOrFile(f.EnvVars, f.FilePath); found {
		if val != "" {
			valInt, err := strconv.ParseUint(val, f.Base, 64)
			if err != nil {
				return fmt.Errorf("could not parse %q as uint value from %s for flag %s: %s", val, source, f.Name, err)
			}

			// Check if the parsed value exceeds the maximum value for uint.
			if valInt > math.MaxUint {
				return fmt.Errorf("value %q from %s for flag %s exceeds maximum uint value", val, source, f.Name)
			}

			f.Value = uint(valInt)
			f.HasBeenSet = true
		}
	}

	for _, name := range f.Names() {
		if f.Destination != nil {
			set.UintVar(f.Destination, name, f.Value, f.Usage)
			continue
		}
		set.Uint(name, f.Value, f.Usage)
	}

	return nil
}

// RunAction executes flag action if set
func (f *UintFlag) RunAction(c *Context) error {
	if f.Action != nil {
		return f.Action(c, c.Uint(f.Name))
	}

	return nil
}

// GetValue returns the flags value as string representation and an empty
// string if the flag takes no value at all.
func (f *UintFlag) GetValue() string {
	return fmt.Sprintf("%d", f.Value)
}

// GetDefaultText returns the default text for this flag
func (f *UintFlag) GetDefaultText() string {
	if f.DefaultText != "" {
		return f.DefaultText
	}
	if f.defaultValueSet {
		return fmt.Sprintf("%d", f.defaultValue)
	}
	return fmt.Sprintf("%d", f.Value)
}

// GetEnvVars returns the env vars for this flag
func (f *UintFlag) GetEnvVars() []string {
	return f.EnvVars
}

// Get returns the flag’s value in the given Context.
func (f *UintFlag) Get(ctx *Context) uint {
	return ctx.Uint(f.Name)
}

// Uint looks up the value of a local UintFlag, returns
// 0 if not found
func (cCtx *Context) Uint(name string) uint {
	if fs := cCtx.lookupFlagSet(name); fs != nil {
		return lookupUint(name, fs)
	}
	return 0
}

func lookupUint(name string, set *flag.FlagSet) uint {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseUint(f.Value.String(), 0, 64)
		if err != nil {
			return 0
		}
		// Check if the parsed value exceeds the maximum value for uint.
		if parsed > math.MaxUint {
			return 0 // Return 0 if the value is out of range for uint.
		}
		return uint(parsed)
	}
	return 0
}
