package utils

// LabelsForPostgreSQL is a function which creates default labels
func LabelsForPostgreSQL(name string) map[string]string {
	return map[string]string{
		"app":                         "postgresql",
		"postgresql_cr":               name,
		"app.kubernetes.io/name":      "postgresql",
		"app.kubernetes.io/instance":  name,
		"app.kubernetes.io/component": "postgresql",
	}
}
