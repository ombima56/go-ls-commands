package listfiles

import (
	"fmt"
	"os"
)

func PrintFileName(file os.FileInfo) {
	color := Reset
	if file.IsDir() {
		color = Blue
	}
	fmt.Printf("%s%s%s ", color, file.Name(), Reset)
}
