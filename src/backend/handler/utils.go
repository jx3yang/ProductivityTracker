package handler

func moveElement(list []string, srcIdx int, destIdx int) []string {
	newList := make([]string, len(list))
	newList[destIdx] = list[srcIdx]

	i := 0

	for idx, element := range list {
		if i == destIdx {
			i += 1
		}
		if idx == srcIdx {
			continue
		}

		newList[i] = element
		i += 1
	}

	return newList
}

func removeOneFromList(list []string, idx int) []string {
	var newList []string

	for i, element := range list {
		if i != idx {
			newList = append(newList, element)
		}
	}

	return newList
}

func addOneToList(list []string, idx int, elem string) []string {
	var newList []string

	if len(list) == 0 {
		newList = append(newList, elem)
		return newList
	}

	for i, element := range list {
		if i == idx {
			newList = append(newList, elem)
		}
		newList = append(newList, element)
	}

	if idx == len(list) {
		newList = append(newList, elem)
	}

	return newList
}

func findFirstIndex(list []string, isMatch func(elem string) bool) int {
	for idx, elem := range list {
		if isMatch(elem) {
			return idx
		}
	}
	return -1
}
