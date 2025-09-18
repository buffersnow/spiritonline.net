package util

import (
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func ToUTF16(endianness unicode.Endianness, input string) ([]byte, error) {
	encoder := unicode.UTF16(endianness, unicode.IgnoreBOM).NewEncoder()
	bytes, _, err := transform.Bytes(encoder, []byte(input))
	if err != nil {
		return nil, fmt.Errorf("util: x/text: %w", err)
	}
	return bytes, nil
}

func FromUTF16(endianness unicode.Endianness, input []byte) (string, error) {
	decoder := unicode.UTF16(endianness, unicode.IgnoreBOM).NewDecoder()
	decoded, _, err := transform.Bytes(decoder, input)
	if err != nil {
		return "", fmt.Errorf("util: x/text: %w", err)
	}
	return string(decoded), nil
}

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
	// expand slices into flat args
	flatArgs := make([]any, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		case []string:
			for _, s := range v {
				flatArgs = append(flatArgs, s)
			}
		case []int64:
			for _, n := range v {
				flatArgs = append(flatArgs, n)
			}
		case []any:
			flatArgs = append(flatArgs, v...)
		default:
			flatArgs = append(flatArgs, v)
		}
	}

	var b strings.Builder
	argIndex := 0

	for i := 0; i < len(query); i++ {
		if query[i] == '?' && argIndex < len(flatArgs) {
			b.WriteString(quoteArg(flatArgs[argIndex]))
			argIndex++
		} else {
			b.WriteByte(query[i])
		}
	}

	// collapse whitespace and trim
	return strings.Join(strings.Fields(b.String()), " ")
}

func quoteArg(arg any) string {
	switch v := arg.(type) {
	case string:
		// escape single quotes
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case []byte:
		// hex representation
		return fmt.Sprintf("'%x'", v)
	case nil:
		return "NULL"
	case []string:
		// return only the first element quoted
		// (the caller should pass elements one by one for multiple ? placeholders)
		if len(v) == 0 {
			return "NULL"
		}
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v[0], "'", "''"))
	default:
		return fmt.Sprintf("'%v'", v)
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

func WithinTimezoneDrift(t time.Time) bool {
	_, offsetSeconds := t.Zone()
	offsetHours := offsetSeconds / 3600

	// Valid timezone drift is between -12 and +14 hours
	return offsetHours >= -12 && offsetHours <= 14
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
