package utils

// NewLabels is a function which creates default labels
func NewLabels(name string, instance string, component string) map[string]string {
	return map[string]string{
		"app":                         "matriarch",
		"postgresql_cr":               instance,
		"app.kubernetes.io/name":      name,
		"app.kubernetes.io/instance":  instance,
		"app.kubernetes.io/component": component,
	}
}
