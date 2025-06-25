package util

import (
	"math/rand"
)

func RandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Batch(funcs []func() error) error {
	for _, fn := range funcs {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
