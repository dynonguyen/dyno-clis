package utils

import (
	"flag"
	"fmt"
	"strings"
)

type FlagItem struct {
	Name, Desc  string
	Required    bool
	Flags       []string
	StrVal      *string
	BoolVal     *bool
	IntVal      *int
	SliceStrVal []string
	DefaultVal  any
	Example     string
}

func (fi *FlagItem) parseFlag() {
	for _, flagName := range fi.Flags {
		switch {
		case fi.StrVal != nil:
			dVal := ""
			if fi.DefaultVal != nil {
				dVal = fi.DefaultVal.(string)
			}
			flag.StringVar(fi.StrVal, flagName, dVal, fi.Desc)
		case fi.BoolVal != nil:
			dVal := false
			if fi.DefaultVal != nil {
				dVal = fi.DefaultVal.(bool)
			}
			flag.BoolVar(fi.BoolVal, flagName, dVal, fi.Desc)
		case fi.IntVal != nil:
			dVal := 0
			if fi.DefaultVal != nil {
				dVal = fi.DefaultVal.(int)
			}
			flag.IntVar(fi.IntVal, flagName, dVal, fi.Desc)
		}
	}
}

func ParseFlags(items []FlagItem, example string) {
	for _, item := range items {
		item.parseFlag()
	}

	flag.Usage = func() {
		fmt.Println("Usage:")

		for _, item := range items {
			keys := make([]string, len(item.Flags))
			for i, f := range item.Flags {
				if len(f) == 1 {
					keys[i] = "-" + f
				} else {
					keys[i] = "--" + f
				}
			}

			flagExample := ""
			if item.Example != "" {
				flagExample = fmt.Sprintf(" (example: %s)", item.Example)
			}

			defaultVal := ""
			if item.DefaultVal != nil && item.DefaultVal != "" {
				defaultVal = fmt.Sprintf(" (default: %v)", item.DefaultVal)
			}

			required := ""
			if item.Required {
				required = " (required)"
			}

			fmt.Printf("  %s: %s%s%s%s\n", strings.Join(keys, ", "), item.Desc, defaultVal, required, flagExample)
		}

		if example != "" {
			fmt.Printf("\nExample:\n  %s\n", example)
		}
	}

	flag.Parse()
}
