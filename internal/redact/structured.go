package redact

import (
	"bytes"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"

	"github.com/mrgoonie/maskara/internal/detect"
)

func redactStructured(path string, original []byte) ([]byte, int, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return redactJSON(original)
	case ".jsonl":
		return redactJSONL(original)
	default:
		return nil, 0, errors.New("structured redaction unsupported for file type")
	}
}

func redactJSON(original []byte) ([]byte, int, error) {
	var value any
	decoder := json.NewDecoder(bytes.NewReader(original))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return nil, 0, err
	}
	replaced := redactJSONValue(&value)
	if replaced == 0 {
		return original, 0, nil
	}
	rewritten, err := json.Marshal(value)
	if err != nil {
		return nil, 0, err
	}
	if bytes.HasSuffix(original, []byte{'\n'}) {
		rewritten = append(rewritten, '\n')
	}
	return rewritten, replaced, nil
}

func redactJSONL(original []byte) ([]byte, int, error) {
	lines := bytes.SplitAfter(original, []byte{'\n'})
	rewritten := make([]byte, 0, len(original))
	replaced := 0

	for _, rawLine := range lines {
		if len(rawLine) == 0 {
			continue
		}
		lineEnding := lineEnding(rawLine)
		line := bytes.TrimSuffix(rawLine, []byte{'\n'})
		line = bytes.TrimSuffix(line, []byte{'\r'})
		if len(bytes.TrimSpace(line)) == 0 {
			rewritten = append(rewritten, rawLine...)
			continue
		}

		var value any
		decoder := json.NewDecoder(bytes.NewReader(line))
		decoder.UseNumber()
		if err := decoder.Decode(&value); err != nil {
			return nil, 0, err
		}
		lineReplaced := redactJSONValue(&value)
		if lineReplaced == 0 {
			rewritten = append(rewritten, rawLine...)
		} else {
			encoded, err := json.Marshal(value)
			if err != nil {
				return nil, 0, err
			}
			rewritten = append(rewritten, encoded...)
			rewritten = append(rewritten, lineEnding...)
		}
		replaced += lineReplaced
	}
	return rewritten, replaced, nil
}

func lineEnding(line []byte) []byte {
	switch {
	case bytes.HasSuffix(line, []byte{'\r', '\n'}):
		return []byte{'\r', '\n'}
	case bytes.HasSuffix(line, []byte{'\n'}):
		return []byte{'\n'}
	default:
		return nil
	}
}

func redactJSONValue(value *any) int {
	switch typed := (*value).(type) {
	case string:
		rewritten, replaced := redactString(typed)
		if replaced > 0 {
			*value = rewritten
		}
		return replaced
	case []any:
		replaced := 0
		for index := range typed {
			replaced += redactJSONValue(&typed[index])
		}
		return replaced
	case map[string]any:
		replaced := 0
		for key, nested := range typed {
			replaced += redactJSONValue(&nested)
			typed[key] = nested
		}
		return replaced
	default:
		return 0
	}
}

func redactString(value string) (string, int) {
	findings := detect.Find(value, "", "", detect.DefaultRules())
	rewritten, replaced := applyRawRedactions([]byte(value), findings)
	return string(rewritten), replaced
}
