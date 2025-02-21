package util

func Modify[T any](c T, f func(c *T)) T {
	f(&c)
	return c
}
