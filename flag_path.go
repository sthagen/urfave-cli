package cli

import (
	"flag"
	"fmt"
)

type Path = string

// GetValue returns the flags value as string representation and an empty
// string if the flag takes no value at all.
func (f *PathFlag) GetValue() string {
	return f.Value
}

func (f *PathFlag) getValueAsAny() (any, error) {
	return f.Value, nil
}

// GetDefaultText returns the default text for this flag
func (f *PathFlag) GetDefaultText() string {
	if f.DefaultText != "" {
		return f.DefaultText
	}
	if f.defaultValue == "" {
		return f.defaultValue
	}
	return fmt.Sprintf("%q", f.defaultValue)
}

// Apply populates the flag given the flag set and environment
func (f *PathFlag) Apply(set *flag.FlagSet) error {
	// set default value so that environment wont be able to overwrite it
	f.defaultValue = f.Value

	if val, _, found := flagFromEnvOrFile(f.EnvVars, f.FilePath); found {
		f.Value = val
		f.HasBeenSet = true
	}

	for _, name := range f.Names() {
		if f.Destination != nil {
			set.StringVar(f.Destination, name, f.Value, f.Usage)
			continue
		}
		set.String(name, f.Value, f.Usage)
	}

	return nil
}

// Get returns the flag’s value in the given Context.
func (f *PathFlag) Get(ctx *Context) string {
	return ctx.Path(f.Name)
}

// RunAction executes flag action if set
func (f *PathFlag) RunAction(c *Context) error {
	if f.Action != nil {
		return f.Action(c, c.Path(f.Name))
	}

	return nil
}

// Path looks up the value of a local PathFlag, returns
// "" if not found
func (cCtx *Context) Path(name string) string {
	if _, flCfg := cCtx.lookupFlagSet(name); flCfg != nil {
		if v, err := cCtx.lookupValue(
			flCfg,
			name,
			func(s string) (any, error) { return s, nil },
		); err == nil {
			return v.(string)
		}
	}

	return ""
}
