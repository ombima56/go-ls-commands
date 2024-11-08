package sorting

func SortFiles(paths []string) {
	n := len(paths)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			// Compare and prioritize lowercase-first
			if shouldSwap(paths[j], paths[j+1]) {
				// Swap
				paths[j], paths[j+1] = paths[j+1], paths[j]
			}
		}
	}
}
