package protocol

import (
	"errors"
	"fmt"
	"net/url"
)

func EnsureScheme(addr string) (string, error) {
	if addr == "" {
		return addr, errors.New("empty addr")
	}

	u, err := url.Parse(addr)
	if err != nil {
		return addr, fmt.Errorf("url: %w", err)
	}

	if u.Scheme == "" {
		return "http://" + u.String(), nil
	}

	return u.String(), err
}
