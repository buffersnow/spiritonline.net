package net

import "time"

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
