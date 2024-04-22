package core

import "fmt"

func InterfaceSliceToStringSlice(data []interface{}) []string {
	strings := make([]string, len(data))
	for i, v := range data {
		// Check if the element is already a string
		if str, ok := v.(string); ok {
			strings[i] = str
		} else {
			// If it's not a string, convert it to a string using fmt.Sprintf or other appropriate method
			strings[i] = fmt.Sprintf("%v", v)
		}
	}
	return strings
}
