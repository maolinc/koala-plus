package collectx

func Contains[T int64 | int32 | string](list []T, item T) bool {
	for i := range list {
		if list[i] == item {
			return true
		}
	}
	return false
}
