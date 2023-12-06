package minion

func channelBufferFull[T any](ch chan T) bool {
	return len(ch) == cap(ch)
}

func channelBufferRemaining[T any](ch chan T) int {
	return cap(ch) - len(ch)
}
