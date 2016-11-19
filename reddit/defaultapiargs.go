package reddit

var (
	defaultAPIArgs = map[string]string{
		"raw_json": "1",
	}
)

// withDefaultAPIArgs returns a value map with the default arguments which
// should be sent to be Reddit set in it.
func withDefaultAPIArgs(m map[string]string) map[string]string {
	if m == nil {
		m = make(map[string]string)
	}
	for k, v := range defaultAPIArgs {
		m[k] = v
	}
	return m
}
