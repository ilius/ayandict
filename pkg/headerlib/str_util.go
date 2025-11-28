package headerlib

func joinWithMaxLen(strs []string, sep string, maxLen int) ([]string, string) {
	if len(strs) == 0 {
		return strs, ""
	}
	if len(strs) == 1 {
		return strs, strs[0]
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
	return strs[:index], out
}
