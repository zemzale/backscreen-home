package slices

// FilterInPlace removes all elements from s for which f returns false.
func FilterInPlace[E any](s []E, f func(E) bool) []E {
	var i int
	for _, v := range s {
		if f(v) {
			s[i] = v
			i++
		}
	}
	return s[:i]
}
