package slices

func Map[S ~[]E, E, R any](s S, f func(E) R) []R {
	r := make([]R, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}
