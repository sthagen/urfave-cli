package cli

import (
	"flag"
	"fmt"
	"strconv"
)

// Apply populates the flag given the flag set and environment
func (f *UintFlag) Apply(set *flag.FlagSet) error {
	// set default value so that environment wont be able to overwrite it
	f.defaultValue = f.Value

	if val, source, found := flagFromEnvOrFile(f.EnvVars, f.FilePath); found {
		if val != "" {
			valInt, err := strconv.ParseUint(val, f.Base, 64)
			if err != nil {
				return fmt.Errorf("could not parse %q as uint value from %s for flag %s: %s", val, source, f.Name, err)
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
	return fmt.Sprintf("%d", f.defaultValue)
}

// Get returns the flag’s value in the given Context.
func (f *UintFlag) Get(ctx *Context) uint {
	return ctx.Uint(f.Name)
}

// RunAction executes flag action if set
func (f *UintFlag) RunAction(c *Context) error {
	if f.Action != nil {
		return f.Action(c, c.Uint(f.Name))
	}

	return nil
}

// Uint looks up the value of a local UintFlag, returns
// 0 if not found
func (cCtx *Context) Uint(name string) uint {
	if _, flCfg := cCtx.lookupFlagSet(name); flCfg != nil {
		parsed, err := strconv.ParseUint(flCfg.Value(), 0, 64)
		if err != nil {
			return 0
		}

		return uint(parsed)
	}

	return 0
}
