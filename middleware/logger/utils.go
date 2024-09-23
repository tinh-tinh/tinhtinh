package logger

import "regexp"

func getLevelName(level Level) string {
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func extractAllContent(s string) []string {
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
