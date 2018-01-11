package compare

import (
	"fmt"
	"reflect"
)

// Kind regards  https://github.com/hashicorp/terraform/blob/master/flatmap/flatten.go for implementation details
func Flatten(thing map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for k, raw := range thing {
		flatten(result, k, reflect.ValueOf(raw))
	}

	return result
}

func flatten(result map[string]string, prefix string, v reflect.Value) {
	v = value(v)

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			result[prefix] = "true"
		} else {
			result[prefix] = "false"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result[prefix] = fmt.Sprintf("%d", v.Int())
	case reflect.Float32, reflect.Float64:
		result[prefix] = fmt.Sprintf("%d", v.Float())
	case reflect.Map:
		flattenMap(result, prefix, v)
	case reflect.Slice:
		flattenSlice(result, prefix, v)
	case reflect.String:
		result[prefix] = v.String()
	case reflect.Invalid:
		fmt.Println(v)
		fmt.Println("invalid")
	case reflect.Ptr:
		flatten(result, prefix, reflect.Indirect(v))

	default:
		fmt.Println(v.Kind())
		fmt.Println(reflect.TypeOf(v))
		fmt.Println(reflect.TypeOf(v).Kind())
		panic(fmt.Sprintf("Unknown: %s", v))
	}
}

func flattenMap(result map[string]string, prefix string, v reflect.Value) {
	for _, k := range v.MapKeys() {
		k = value(k)

		if k.Kind() != reflect.String {
			panic(fmt.Sprintf("%s: map key is not string: %s", prefix, k))
		}

		flatten(result, fmt.Sprintf("%s.%s", prefix, k.String()), v.MapIndex(k))
	}
}

func flattenSlice(result map[string]string, prefix string, v reflect.Value) {
	prefix = prefix + "."

	result[prefix+"#"] = fmt.Sprintf("%d", v.Len())
	for i := 0; i < v.Len(); i++ {
		flatten(result, fmt.Sprintf("%s%d", prefix, i), v.Index(i))
	}
}

func value(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}
