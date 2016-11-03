package api

var (
	defaultValues = map[string]string{
		"raw_json": "1",
	}
)

// withDefaults returns a value map with the defaults set in it.
func withDefaults(m map[string]string) map[string]string {
	if m == nil {
		m = make(map[string]string)
	}
	for k, v := range defaultValues {
		m[k] = v
	}
	return m
}
