// internal/command/flags.go
package command

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Flag represents a command flag
type Flag struct {
	Name      string
	Shorthand string
	Usage     string
	Type      string // "bool", "string", "int"
	Default   interface{}
}

// FlagDefiner is an optional interface for commands that define flags
type FlagDefiner interface {
	DefineFlags() []Flag
}

// SetupFlags adds flags to a cobra command
func SetupFlags(cobraCmd *cobra.Command, flags []Flag) {
	for _, flag := range flags {
		switch flag.Type {
		case "bool":
			defaultVal, _ := flag.Default.(bool)
			if flag.Shorthand != "" {
				cobraCmd.Flags().BoolP(flag.Name, flag.Shorthand, defaultVal, flag.Usage)
			} else {
				cobraCmd.Flags().Bool(flag.Name, defaultVal, flag.Usage)
			}
		case "string":
			defaultVal, _ := flag.Default.(string)
			if flag.Shorthand != "" {
				cobraCmd.Flags().StringP(flag.Name, flag.Shorthand, defaultVal, flag.Usage)
			} else {
				cobraCmd.Flags().String(flag.Name, defaultVal, flag.Usage)
			}
		case "int":
			defaultVal, _ := flag.Default.(int)
			if flag.Shorthand != "" {
				cobraCmd.Flags().IntP(flag.Name, flag.Shorthand, defaultVal, flag.Usage)
			} else {
				cobraCmd.Flags().Int(flag.Name, defaultVal, flag.Usage)
			}
		}
	}
}

// ParseFlags extracts flag values from cobra command
func ParseFlags(cobraCmd *cobra.Command) map[string]interface{} {
	flags := make(map[string]interface{})

	cobraCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			switch f.Value.Type() {
			case "bool":
				val, _ := cobraCmd.Flags().GetBool(f.Name)
				flags[f.Name] = val
			case "string":
				val, _ := cobraCmd.Flags().GetString(f.Name)
				flags[f.Name] = val
			case "int":
				val, _ := cobraCmd.Flags().GetInt(f.Name)
				flags[f.Name] = val
			default:
				flags[f.Name] = f.Value.String()
			}
		}
	})

	return flags
}
