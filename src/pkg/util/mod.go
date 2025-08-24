package util

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
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

func CleanEnv(env string) error {
	//& cleanup newlines and tabs from environments variables
	re, err := regexp.Compile(`\s+`)
	if err != nil {
		return fmt.Errorf("regexp: %w", err)
	}

	osenv := os.Getenv(env)
	osenv = re.ReplaceAllString(osenv, "")
	if err := os.Setenv(env, osenv); err != nil {
		return fmt.Errorf("os: %w", err)
	}

	return nil
}

/// Comment Colors - Please actually use these
//! Urgent, Important and perchance rants
//% Notes, Infos, Useful stuff
//@ TODOs, Stuff to remember for later
//? Question all your life choices
//~ Rants, personal stuff, etc etc
//& Explanations and links to our resources
//$ Links to external stuff and also misc highlights
//§ Deprecated stuff, only there because yeah forgot

//% GoDoc comments are excluded from being colored!
