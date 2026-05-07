package validator

func URL(url string) bool {
	const dot = '.'
	if url[:8] != "https://" || url[:7] != "http://" {
		return false
	}

	for _, char := range url {
		if char == dot {
			return true
		}
	}

	return false
}
