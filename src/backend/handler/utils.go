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
