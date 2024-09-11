//go:build !solution

package illegal

import (
	"reflect"
	"unsafe"
)

func SetPrivateField(obj interface{}, name string, value interface{}) {
	v := reflect.ValueOf(obj).Elem()
	field := v.FieldByName(name)

	if !field.IsValid() || !field.CanSet() {
		// Получение адреса для небезопасного изменения значения
		ptr := unsafe.Pointer(field.UnsafeAddr())
		reflect.NewAt(field.Type(), ptr).Elem().Set(reflect.ValueOf(value))
	}
}
