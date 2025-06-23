package util

func CastMessage[T any](body any) T {
	return body.(T)
}
