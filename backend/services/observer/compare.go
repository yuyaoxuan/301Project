package observer

import (
	"reflect"
)

// Compare compares two objects and returns the modified fields
func Compare(before, after interface{}) map[string]interface{} {
	changes := make(map[string]interface{})

	// Use reflection to compare fields dynamically
	beforeVal := reflect.ValueOf(before).Elem()
	afterVal := reflect.ValueOf(after).Elem()

	for i := 0; i < beforeVal.NumField(); i++ {
		beforeField := beforeVal.Field(i)
		afterField := afterVal.Field(i)

		// Compare if the values are different
		if !reflect.DeepEqual(beforeField.Interface(), afterField.Interface()) {
			fieldName := beforeVal.Type().Field(i).Name
			changes[fieldName] = map[string]interface{}{
				"before": beforeField.Interface(),
				"after":  afterField.Interface(),
			}
		}
	}

	return changes
}
