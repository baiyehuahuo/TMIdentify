package utils

func Deduplication(slice []string) []string {
	set := make(map[string]struct{})
	for i := range slice {
		set[slice[i]] = struct{}{}
	}
	result := make([]string, 0, len(set))
	for key := range set {
		result = append(result, key)
	}
	return result
}
