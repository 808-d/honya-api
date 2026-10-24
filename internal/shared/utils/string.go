package utils

// OptionalString returns nil for an empty string, otherwise a pointer to it.
func OptionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
