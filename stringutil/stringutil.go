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

func ArrayIndex(texts []string, text string) int {
	if texts != nil && len(texts) > 0 {
		for i, t := range texts {
			if t == text {
				return i
			}
		}
	}
	return -1
}

func ArrayContains(texts []string, text string) bool {
	return ArrayIndex(texts, text) != -1
}

func ArrayContainsAll(container, elements []string) bool {
	if elements != nil && len(elements) > 0 {
		for _, elem := range elements {
			if !ArrayContains(container, elem) {
				return false
			}
		}
	}
	return true
}

func SetEqual(a, b []string) bool {
	return ArrayContainsAll(a, b) && ArrayContainsAll(b, a)
}
