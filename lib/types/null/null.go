package null

import "reflect"

// All returns true if all given values are nil.
func All(values ...any) bool {
	for _, v := range values {
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
			if rv.IsNil() {
				continue
			}
			return false
		case reflect.Invalid:
			continue
		default:
			return false
		}
	}
	return true
}
