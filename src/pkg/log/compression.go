package log

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func (l Logger) archiveLogJob() {
	for {
		time.Sleep(24 * time.Hour)
		if err := l.archiveLog(); err != nil {
			l.Error("Log Archiver", "Unable to complete daily log archival job: %v", err)
		}
	}
}

func (l *Logger) archiveLog() error {

	// Check if the file exists
	if _, err := os.Stat(l.fileName); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("log: %w", err)
		}

		l.Warning("Log Archiver", "No log file for archiving found!")
		return nil
	}

	l.mu.Lock()

	logFileRWHandle, err := os.OpenFile(l.fileName, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("log: %w", err)
	}
	defer logFileRWHandle.Close()

	date := time.Now().Local().Format("02-01-2006")
	zipFile, err := os.OpenFile(
		fmt.Sprintf("logs-%s.zip", date),
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

	l.mu.Unlock()
	l.Event("Log Archiver", "Archived logs for %s", date)

	return nil
}
