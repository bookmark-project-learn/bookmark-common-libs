package array

func SplitIntoBatch[T any](object []T, batchSize int) [][]T {
	if batchSize <= 0 {
		return [][]T{}
	}
	batches := make([][]T, 0, (len(object)+batchSize-1)/batchSize)
	for i := 0; i < len(object); i += batchSize {
		end := i + batchSize
		if end > len(object) {
			end = len(object)
		}
		batches = append(batches, object[i:end])
	}
	return batches
}
