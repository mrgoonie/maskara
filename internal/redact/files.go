package redact

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func backupName(path string) string {
	base := path + ".maskara.bak"
	if _, err := os.Stat(base); errors.Is(err, os.ErrNotExist) {
		return base
	}
	return base + "." + time.Now().UTC().Format("20060102150405")
}

func writeReplace(path string, content []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	temp, err := os.CreateTemp(dir, filepath.Base(path)+".maskara-*")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tempPath)
		}
	}()
	if _, err := temp.Write(content); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Chmod(mode); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tempPath, path); err != nil {
		if runtime.GOOS != "windows" {
			return err
		}
		if removeErr := os.Remove(path); removeErr != nil {
			return err
		}
		if renameErr := os.Rename(tempPath, path); renameErr != nil {
			return renameErr
		}
	}
	cleanup = false
	return nil
}

func validateStructured(path string, before, after []byte) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		trimmed := bytes.TrimSpace(before)
		if len(trimmed) > 0 && json.Valid(trimmed) && !json.Valid(bytes.TrimSpace(after)) {
			return errors.New("redaction would break JSON: " + path)
		}
	case ".jsonl":
		if jsonlValid(before) && !jsonlValid(after) {
			return errors.New("redaction would break JSONL: " + path)
		}
		if countLines(before) != countLines(after) {
			return errors.New("redaction would change JSONL line count: " + path)
		}
	}
	return nil
}

func jsonlValid(content []byte) bool {
	lines := bytes.Split(content, []byte{'\n'})
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if !json.Valid(line) {
			return false
		}
	}
	return true
}

func countLines(content []byte) int {
	if len(content) == 0 {
		return 0
	}
	return bytes.Count(content, []byte{'\n'}) + 1
}
