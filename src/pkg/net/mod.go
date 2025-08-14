package net

import (
	"fmt"
	"time"

	"buffersnow.com/spiritonline/pkg/log"
)

type NetUtils struct {
	log *log.Logger
}

func New(log *log.Logger) (*NetUtils, error) {
	return &NetUtils{log: log}, nil
}

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
