package common

func CloneMap(orignal map[string]string) map[string]string {
	cloned := make(map[string]string, len(orignal))
	for k, v := range orignal {
		cloned[k] = v
	}
	return cloned
}

func MergeMaps(original, toMerge map[string]string) {
	for key, value := range toMerge {
		original[key] = value // Adds or updates key-value pairs
	}
}
