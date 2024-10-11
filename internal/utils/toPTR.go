package utils

func ToPTR[t any](v t) *t {
	return &v
}
