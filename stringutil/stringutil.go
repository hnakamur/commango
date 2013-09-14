package stringutil

func Uniq(lines []string) []string {
	lineSet := make(map[string]int)

	i := 0
	for _, line := range(lines) {
		if _, ok := lineSet[line]; !ok {
			lineSet[line] = i
			i++
		}
	}

	result := make([]string, len(lineSet))
	for line, i := range(lineSet) {
		result[i] = line
	}
	return result
}
