package logger

import (
	"fmt"
	"os"
	"regexp"
)

func GetLevelName(level Level) string {
	switch level {
	case LevelFatal:
		return "fatal"
	case LevelError:
		return "error"
	case LevelWarn:
		return "warn"
	case LevelInfo:
		return "info"
	case LevelDebug:
		return "debug"
	default:
		return ""
	}
}

func ExtractAllContent(s string) []string {
	re := regexp.MustCompile(`\$\{(.*?)\}`)
	matches := re.FindAllStringSubmatch(s, -1)

	var results []string
	for _, match := range matches {
		if len(match) > 1 {
			results = append(results, match[1])
		}
	}

	return results
}

func checkAvailableFile(filename string, max int64) bool {
	if max <= 0 {
		return true
	}
	fi, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// no file yet = available
		return true
	}
	if err != nil {
		fmt.Printf("Failed to check log file: %v\n", err)
		return false
	}
	return fi.Size() < max*MiB
}

func appendMetadata(base Metadata, extra ...Metadata) Metadata {
	merged := make(Metadata)
	for k, v := range base {
		merged[k] = v
	}
	if len(extra) > 0 {
		for _, m := range extra {
			for k, v := range m {
				merged[k] = v
			}
		}
	}

	return merged
}
