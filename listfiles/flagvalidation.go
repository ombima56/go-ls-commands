package listfiles

import (
	"fmt"
	"strings"
)

func ValidateFlags(args []string) (bool, bool, bool, bool, bool, error) {
	var longFlag, allFlag, recursiveFlag, timeFlag, reverseFlag bool

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// Handle combined flags (-la)
			flagStr := strings.TrimPrefix(arg, "-")
			if strings.HasPrefix(flagStr, "-") {
				// Handle long flags (--long)
				flagStr = strings.TrimPrefix(flagStr, "-")
				switch flagStr {
				case "long":
					longFlag = true
				case "all":
					allFlag = true
				case "recursive":
					recursiveFlag = true
				case "time":
					timeFlag = true
				case "reverse":
					reverseFlag = true
				default:
					return false, false, false, false, false, fmt.Errorf("invalid option --%s", flagStr)
				}
			} else {
				// Handle short flags (-l)
				for _, flag := range flagStr {
					switch flag {
					case 'l':
						longFlag = true
					case 'a':
						allFlag = true
					case 'R':
						recursiveFlag = true
					case 't':
						timeFlag = true
					case 'r':
						reverseFlag = true
					default:
						return false, false, false, false, false, fmt.Errorf("invalid option -- '%c'", flag)
					}
				}
			}
		}
	}

	return longFlag, allFlag, recursiveFlag, timeFlag, reverseFlag, nil
}
