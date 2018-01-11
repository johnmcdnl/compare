package compare

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Fill(i interface{}) {
	setData(reflect.ValueOf(i))
}

func setData(v reflect.Value) {

	if v.Kind() != reflect.Ptr {
		fmt.Println(v.Kind())
		panic("v.Kind() != reflect.Ptr")
	}

	v = reflect.Indirect(v)
	switch v.Kind() {

	case reflect.Int:
		v.Set(reflect.ValueOf(int(rand.Intn(100))))
	case reflect.Int8:
		v.Set(reflect.ValueOf(int8(rand.Intn(100))))
	case reflect.Int16:
		v.Set(reflect.ValueOf(int16(rand.Intn(100))))
	case reflect.Int32:
		v.Set(reflect.ValueOf(int32(rand.Intn(100))))
	case reflect.Int64:
		v.Set(reflect.ValueOf(int64(rand.Intn(100))))
	case reflect.Float32:
		v.Set(reflect.ValueOf(rand.Float32()))
	case reflect.Float64:
		v.Set(reflect.ValueOf(rand.Float64()))
	case reflect.String:
		v.SetString(randomString(25))
	case reflect.Bool:
		val := rand.Intn(2) > 0
		v.SetBool(val)
	case reflect.Slice:
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			setData(v.Field(i).Addr())
		}
	case reflect.Ptr:
		panic(fmt.Sprintf("Unhanlded Ptr() -- what to do with this one"))
	default:
		panic(fmt.Sprintf("Unhanlded Kind() %s", v.Kind()))
	}

}

func randomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
