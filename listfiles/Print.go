package listfiles

import (
	"fmt"
	"os"

	"go-ls-commands/colors"
)

func PrintFileName(file os.FileInfo) {
	color := colors.GetFileColor(file)
	fmt.Printf("%s%s%s ", color, file.Name(), colors.Reset)
}
