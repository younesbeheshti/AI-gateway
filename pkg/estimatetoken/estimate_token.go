package estimatetoken

func EstimateTokens(text string) int64 {
	tokens := int64(len(text) / 4)
	if tokens == 0 && len(text) > 0 {
		return 1
	}
	return tokens
}
