package props

func Copy(params map[string]string) map[string]string {
	newMap := make(map[string]string, len(params))
	for k, v := range params {
		newMap[k] = v
	}
	return newMap
}
