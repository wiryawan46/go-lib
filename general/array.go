package general

// InArray ...
// Check if variable exists in an array
func InArray(variable interface{}, array interface{}) bool {
	switch key := variable.(type) {
	case string:
		for _, item := range array.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range array.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range array.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}
