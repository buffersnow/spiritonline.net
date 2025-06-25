package net

import (
	"fmt"
	"time"
)

type NetConnection interface {
	AcceptIncoming() error
	Close() error
	GetRemoteAddress() string
	ReadBytes() (data []byte, err error)
	ReadBytesEx(timeout time.Time) (data []byte, err error)
	ReadText() (data string, err error)
	ReadTextEx(timeout time.Time) (data string, err error)
	WriteBytes(data []byte) error
	WriteText(data string) error
}

func prettyBytes(slice []byte) string {
	var hexString string
	for _, b := range slice {
		hexString += fmt.Sprintf("%02x ", b)
	}
	return fmt.Sprintf("[ %s]", hexString)
}
