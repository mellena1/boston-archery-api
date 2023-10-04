package slices

func Map[T any, U any](s []T, f func(T) U) []U {
	if s == nil {
		return nil
	}

	newS := []U{}
	for _, item := range s {
		newS = append(newS, f(item))
	}
	return newS
}
