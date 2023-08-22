package nilcheck

import "reflect"

// AreAllNil returns true if all given values are nil.
//
// AreAllNil(nil, nil) -> true
// AreAllNil(nil, 1) -> false
func AreAllNil(values ...any) bool {
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

// AnyNotNil returns true if any given value is not nil.
//
// AnyNotNil(nil, nil) -> false
// AnyNotNil(nil, 1) -> true
func AnyNotNil(values ...any) bool {
	return !AreAllNil(values...)
}
