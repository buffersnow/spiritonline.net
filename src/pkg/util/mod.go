package util

import (
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func RandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
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

// Replaces '?' in the query with the provided args !ONLY FOR LOGGING!
func FormatSQL(query string, args ...any) string {
	var b strings.Builder
	argIndex := 0
	for i := 0; i < len(query); i++ {
		if query[i] == '?' && argIndex < len(args) {
			// Write the argument in a quoted form
			b.WriteString(quoteArg(args[argIndex]))
			argIndex++
		} else {
			b.WriteByte(query[i])
		}
	}
	return b.String()
}
func quoteArg(arg any) string {
	switch v := arg.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case []byte:
		return fmt.Sprintf("'%x'", v)
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func CountSQLRows(dest any) int64 {
	if dest == nil {
		return 0
	}
	rv := reflect.ValueOf(dest)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice:
		return int64(rv.Len())
	default:
		// struct, map, primitive, etc. — assume at most 1 row
		return 1
	}
}

func HexToByte(hexStr string) byte {
	if len(hexStr) != 2 {
		return 0xFF
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil || len(bytes) != 1 {
		return 0xFF
	}

	return bytes[0]
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
