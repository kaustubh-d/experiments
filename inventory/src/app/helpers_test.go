package app

// helper: compare two slices of EnvNameType ignoring order
func envSlicesEqualUnordered(a, b []EnvNameType) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int)
	for _, v := range a {
		m[string(v)]++
	}
	for _, v := range b {
		s := string(v)
		if m[s] == 0 {
			return false
		}
		m[s]--
	}
	for _, cnt := range m {
		if cnt != 0 {
			return false
		}
	}
	return true
}
