package util

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"buffersnow.com/spiritonline/pkg/version"
)

func Initialize() {
	settingsWg.Add(1)

	buffer := "- Welcome to spiritonline.net -> " + version.GetService() + " v" + version.GetVersion() + "\n"
	buffer += "Build Tag: " + version.GetPartialTag()
	Log.Raw("\033[38;5;61m", "%s", buffer)
}

func AwaitInterrupt() os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	return <-c
}

func RandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func PrettyBytes(slice []byte) string {
	var hexString string
	for _, b := range slice {
		hexString += fmt.Sprintf("%02x ", b)
	}
	return fmt.Sprintf("[ %s]", hexString)
}
