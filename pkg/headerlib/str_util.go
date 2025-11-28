package headerlib

func joinWithMaxLen(strs []string, sep string, maxLen int) (string, int) {
	if len(strs) == 0 {
		return "", 0
	}
	if len(strs) == 1 {
		return strs[0], 1
	}
	out := strs[0]
	index := 1
	for index < len(strs) {
		newOut := out + sep + strs[index]
		if len(newOut) > maxLen {
			break
		}
		out = newOut
		index++
	}
	return out, index
}
