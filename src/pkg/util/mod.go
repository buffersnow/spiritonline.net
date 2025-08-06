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

/// Comment Colors
//! Urgent, Important and perchance rants
//% Notes, Infos, Useful stuff
//@ TODOs, Stuff to remember for later
//? Question all your life choices
//~ Rants, personal stuff, etc etc
//& Explanations and links to our resources
//$ Links to external stuff and also misc highlights
//- Deprecated stuff, only there because yeah forgot
