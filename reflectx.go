package envutils

import "reflect"

// deref 返回对象类型。 如 typ 是指针则返回指针指向的类型
func deref(typ reflect.Type) reflect.Type {

	if typ.Kind() == reflect.Ptr {
		return typ.Elem()
	}
	return typ
}
