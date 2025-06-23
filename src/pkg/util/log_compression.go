package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"time"
)

func (l utilLog) compressLogJob() {
	for {
		time.Sleep(24 * time.Hour)
		l.compressLog()
	}
}

func (l *utilLog) compressLog() {
	l.writeMutex.Lock()
	defer l.writeMutex.Unlock()

	logFileRWHandle, err := os.OpenFile(l.logFileName, os.O_RDWR, 0644)
	if err != nil {
		l.Warning("Log Compression", "Unable to open logfile to compress! (Might be intentional?) %s", err.Error())
		return
	}
	defer logFileRWHandle.Close()

	zipFile, err := os.OpenFile(
		fmt.Sprintf("logs-%s.zip", time.Now().Local().Format("02-01-2006")),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644,
	)

	if err != nil {
		l.Error("Log Compression", "Failed to create log archive zip %s", err.Error())
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileInfo, err := logFileRWHandle.Stat()
	if err != nil {
		l.Error("Log Compression", "Failed to populate file header for archive zip %s", err.Error())
		return
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		l.Error("Log Compression", "Failed to populate file header for archive zip %s", err.Error())
		return
	}
	header.Name = l.logFileName

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		l.Error("Log Compression", "Failed to create file header zip item %s", err.Error())
		return
	}

	if _, err := io.Copy(writer, logFileRWHandle); err != nil {
		l.Error("Log Compression", "Failed to copy file to zip %s", err.Error())
		return
	}
}
