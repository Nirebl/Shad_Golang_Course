//go:build !solution

package reversemap

import (
	"reflect"
)

func ReverseMap(forward interface{}) interface{} {
	t := reflect.TypeOf(forward)
	forwardMap := reflect.ValueOf(forward)

	reversedMap := reflect.MakeMap(reflect.MapOf(t.Elem(), t.Key()))

	for _, k := range forwardMap.MapKeys() {
		v := forwardMap.MapIndex(k)
		reversedMap.SetMapIndex(v, k)
	}
	return reversedMap.Interface()
}
