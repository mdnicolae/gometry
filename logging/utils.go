package logging

// NormalizeAttributes ensures that attributes are not nil by returning an empty map if nil.
func NormalizeAttributes(attributes []map[string]interface{}) map[string]interface{} {
	if len(attributes) > 0 && attributes[0] != nil {
		return attributes[0]
	}
	return make(map[string]interface{})
}
