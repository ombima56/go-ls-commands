package listfiles

import (
	"fmt"
	"strings"
)

// Options struct to hold all command flags
type Options struct {
	LongFormat    bool
	AllFiles      bool
	Recursive     bool
	SortByTime    bool
	ReverseSort   bool
}

func ValidateFlags(args []string) (Options, error) {
	opts := Options{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// Handle combined flags (-la)
			flagStr := strings.TrimPrefix(arg, "-")
			if strings.HasPrefix(flagStr, "-") {
				// Handle long flags (--long)
				flagStr = strings.TrimPrefix(flagStr, "-")
				switch flagStr {
				case "long":
					opts.LongFormat = true
				case "all":
					opts.AllFiles = true
				case "recursive":
					opts.Recursive = true
				case "time":
					opts.SortByTime = true
				case "reverse":
					opts.ReverseSort = true
				default:
					return Options{}, fmt.Errorf("invalid option --%s", flagStr)
				}
			} else {
				// Handle short flags (-l)
				for _, flag := range flagStr {
					switch flag {
					case 'l':
						opts.LongFormat = true
					case 'a':
						opts.AllFiles = true
					case 'R':
						opts.Recursive = true
					case 't':
						opts.SortByTime = true
					case 'r':
						opts.ReverseSort = true
					default:
						return Options{}, fmt.Errorf("invalid option -- '%c'", flag)
					}
				}
			}
		}
	}

	return opts, nil
}