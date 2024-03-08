package main

// map merge function
func merge(m1, m2 map[string]string) map[string]string {
	merged := map[string]string{}
	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		merged[k] = v
	}
	return merged
}
