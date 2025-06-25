package log

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"time"
)

func (l Logger) archiveLogJob() {
	for {
		time.Sleep(24 * time.Hour)
		err := l.archiveLog()
		if err != nil {
			l.Error("Log Archiver", "Unable to complete daily log archival job: %v", err)
		}
	}
}

func (l *Logger) archiveLog() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	logFileRWHandle, err := os.OpenFile(l.fileName, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("log: %w", err)
	}
	defer logFileRWHandle.Close()

	zipFile, err := os.OpenFile(
		fmt.Sprintf("logs-%s.zip", time.Now().Local().Format("02-01-2006")),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644,
	)

	if err != nil {
		return fmt.Errorf("log: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileInfo, err := logFileRWHandle.Stat()
	if err != nil {
		return fmt.Errorf("log: %w", err)
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return fmt.Errorf("log: %w", err)
	}
	header.Name = l.fileName

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("log: %w", err)
	}

	if _, err := io.Copy(writer, logFileRWHandle); err != nil {
		return fmt.Errorf("log: %w", err)
	}

	return nil
}
