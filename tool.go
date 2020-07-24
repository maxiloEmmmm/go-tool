package go_tool

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func Md5(str string) string {
	var m = md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func Uint8sToBytes(src []uint8) []byte {
	var dst []byte
	for _, b := range src {
		dst = append(dst, byte(b))
	}
	return dst
}

func TryInterfacePtr(data interface{}) reflect.Value {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		return v.Elem()
	} else {
		return v
	}
}

func AssetsMapOrS(kind reflect.Kind, msg string) {
	if kind != reflect.Map && kind != reflect.Struct {
		AssetsError(errors.New(msg))
	}
}

func AssetsSlice(kind reflect.Kind, msg string) {
	if kind != reflect.Slice {
		AssetsError(errors.New(msg))
	}
}

func AssetsMap(kind reflect.Kind, msg string) {
	if kind != reflect.Map {
		AssetsError(errors.New(msg))
	}
}

func AssetsPtr(kind reflect.Kind, msg string) {
	if kind != reflect.Ptr {
		AssetsError(errors.New(msg))
	}
}

func StringJoin(ss ...string) string {
	buffer := new(strings.Builder)
	for _, s := range ss {
		buffer.WriteString(s)
	}
	return buffer.String()
}

func MapToInterface(data interface{}) map[interface{}]interface{} {
	v := TryInterfacePtr(data)

	AssetsMap(v.Kind(), "map转换接口错误, 非map接口")
	keys := v.MapKeys()
	dst := make(map[interface{}]interface{}, len(keys))

	for _, kV := range keys {
		dst[kV.Interface()] = v.MapIndex(kV).Interface()
	}

	return dst
}

func MapMap(data interface{}, cb func(interface{}) interface{}) interface{} {
	v := TryInterfacePtr(data)
	AssetsMap(v.Kind(), "map转换接口错误, 非map接口")
	keys := v.MapKeys()
	dst := make(map[interface{}]interface{}, v.Len())
	for _, key := range keys {
		dst[key.Type().Name()] = cb(v.MapIndex(key).Interface())
	}
	return dst
}

func MapKeyFilter(data interface{}, cb func(interface{}) bool) interface{} {
	v := TryInterfacePtr(data)
	AssetsMap(v.Kind(), "map转换接口错误, 非map接口")
	keys := v.MapKeys()
	dst := make(map[interface{}]interface{}, v.Len())
	for _, key := range keys {
		k := key.Interface()
		if cb(k) {
			dst[k] = v.MapIndex(key).Interface()
		}
	}
	return dst
}

func MapKeyFilterToArray(data interface{}, cb func(interface{}) bool) interface{} {
	v := TryInterfacePtr(data)
	AssetsMap(v.Kind(), "map转换接口错误, 非map接口")
	keys := v.MapKeys()
	dst := make([]interface{}, 0, v.Len())
	for _, key := range keys {
		k := key.Interface()
		if cb(k) {
			dst = append(dst, v.MapIndex(key).Interface())
		}
	}
	return dst
}

func MapValueFilter(data interface{}, cb func(interface{}) bool) interface{} {
	v := TryInterfacePtr(data)
	AssetsMap(v.Kind(), "map转换接口错误, 非map接口")
	keys := v.MapKeys()
	dst := make(map[interface{}]interface{}, v.Len())
	for _, key := range keys {
		value := v.MapIndex(key).Interface()
		if cb(value) {
			dst[key.Interface()] = value
		}
	}
	return dst
}

func MapValueFilterToArray(data interface{}, cb func(interface{}) bool) interface{} {
	v := TryInterfacePtr(data)
	AssetsMap(v.Kind(), "map转换接口错误, 非map接口")
	keys := v.MapKeys()
	dst := make([]interface{}, 0, v.Len())
	for _, key := range keys {
		value := v.MapIndex(key).Interface()
		if cb(value) {
			dst = append(dst, value)
		}
	}
	return dst
}

func MapKeys(data interface{}, cb func(interface{}) interface{}) interface{} {
	v := TryInterfacePtr(data)
	AssetsMap(v.Kind(), "map转换接口错误, 非map接口")
	keys := v.MapKeys()
	return ArrayMap(&keys, func(d interface{}) interface{} {
		k := d.(reflect.Value).Interface()
		if cb != nil {
			return cb(k)
		}
		return d.(reflect.Value).Interface()
	})
}

func MapValues(data interface{}, cb func(interface{}) interface{}) interface{} {
	v := TryInterfacePtr(data)
	AssetsMap(v.Kind(), "map转换接口错误, 非map接口")
	keys := v.MapKeys()
	return ArrayMap(&keys, func(d interface{}) interface{} {
		value := v.MapIndex(d.(reflect.Value)).Interface()
		if cb != nil {
			return cb(value)
		}
		return value
	})
}

func AssetMapSet(do bool, m interface{}, key interface{}, value interface{}) {
	if do {
		v := reflect.ValueOf(m)

		if v.Kind() != reflect.Ptr {
			AssetsSlice(v.Kind(), "非指针, 无法set map")
		} else {
			obj := v.Elem()
			obj.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
		}
	}
}

func ArrayMap(data interface{}, cb func(interface{}) interface{}) interface{} {
	v := TryInterfacePtr(data)
	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	vl := v.Len()
	dst := make([]interface{}, vl)
	for i := 0; i < vl; i++ {
		dst[i] = cb(v.Index(i).Interface())
	}
	return dst
}

func ArrayFilter(data interface{}, cb func(interface{}) bool) interface{} {
	v := TryInterfacePtr(data)
	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	vl := v.Len()
	dst := make([]interface{}, 0, vl)
	for i := 0; i < vl; i++ {
		tmp := v.Index(i).Interface()
		if cb(tmp) {
			dst = append(dst, tmp)
		}
	}
	return dst
}

func ArrayReduce(data interface{}, cb func(float64, interface{}) float64, start float64) float64 {
	v := TryInterfacePtr(data)
	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	vl := v.Len()
	for i := 0; i < vl; i++ {
		start = cb(start, v.Index(i).Interface())
	}
	return start
}

type ArrayKeyByS map[interface{}][]interface{}

func ArrayKeyBy(data interface{}, key string) ArrayKeyByS {
	v := TryInterfacePtr(data)
	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	vl := v.Len()
	dst := make(ArrayKeyByS, vl)
	for i := 0; i < vl; i++ {
		tmp := v.Index(i).Interface()
		// 相信宝贝你不会存在不存在的情况 - -
		k, _ := Get(tmp, key)
		dst[k] = append(dst[k], tmp)
	}
	return dst
}

func ArrayKeyByFunc(data interface{}, key string, cb func(interface{}, interface{}) interface{}) map[interface{}]interface{} {
	v := TryInterfacePtr(data)
	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	vl := v.Len()
	dst := make(map[interface{}]interface{}, vl)

	for i := 0; i < vl; i++ {
		tmp := v.Index(i).Interface()
		// 相信宝贝你不会存在不存在的情况 - -
		k, _ := Get(tmp, key)
		dst[k] = cb(dst[k], tmp)
	}
	return dst
}

func ArrayMakeKey(data interface{}, key string) map[interface{}]interface{} {
	v := TryInterfacePtr(data)
	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	vl := v.Len()
	dst := make(map[interface{}]interface{}, vl)
	for i := 0; i < vl; i++ {
		tmp := v.Index(i).Interface()
		// 相信宝贝你不会存在不存在的情况 - -
		k, _ := Get(tmp, key)
		dst[k] = tmp
	}
	return dst
}

func ArrayMakeKeyFunc(data interface{}, key string, cb func(d interface{}) interface{}) map[interface{}]interface{} {
	v := TryInterfacePtr(data)
	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	vl := v.Len()
	dst := make(map[interface{}]interface{}, vl)
	for i := 0; i < vl; i++ {
		tmp := v.Index(i).Interface()
		// 相信宝贝你不会存在不存在的情况 - -
		k, _ := Get(tmp, key)
		dst[k] = cb(tmp)
	}
	return dst
}

func ArrayFirst(data interface{}) interface{} {
	v := TryInterfacePtr(data)

	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	if v.Len() == 0 {
		AssetsError(errors.New("数组为空, 无法识别子元素类型"))
	}

	return v.Index(0).Interface()
}

func GetArrayType(data interface{}) reflect.Type {
	return reflect.TypeOf(ArrayFirst(data))
}

func GetTypeFieldNum(t reflect.Type) int {
	kind := t.Kind()
	AssetsMapOrS(kind, "key提取接口错误, 非可提取接口")

	if kind == reflect.Map {
		return t.Len()
	} else {
		return t.NumField()
	}
}

func GetTypeFieldBySet(data interface{}, keys []string) []string {
	v := reflect.ValueOf(data)
	kind := v.Kind()
	AssetsMapOrS(kind, "无法提取字段类型")

	fieldNum := GetTypeFieldNum(v.Type())
	fields := make([]string, 0, fieldNum)
	if kind == reflect.Map {
		fields = ArrayMap(v.MapKeys(), func(d interface{}) interface{} {
			return d.(reflect.Value).Interface()
		}).([]string)
	} else {
		t := reflect.TypeOf(data)
		for i := 0; i < fieldNum; i++ {
			fields = append(fields, t.Field(i).Name)
		}
	}

	dstKeys := make([]string, 0, fieldNum)
	for _, k := range fields {
		if !InArray(keys, k) {
			dstKeys = append(dstKeys, k)
		}
	}
	return dstKeys
}

func InArray(data interface{}, find interface{}) bool {
	v := TryInterfacePtr(data)
	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	dNum := v.Len()
	exist := false
	for i := 0; i < dNum; i++ {
		exist = v.Index(i).Interface() == find
		if exist {
			break
		}
	}
	return exist
}

type ArrayPickS []map[string]interface{}

func ArrayPick(data interface{}, keys []string) ArrayPickS {
	v := TryInterfacePtr(data)
	AssetsSlice(v.Kind(), "数组转换接口错误, 非数组接口")

	dNum := v.Len()
	dst := make(ArrayPickS, dNum)
	dstT := GetArrayType(data).Kind()
	fieldNum := len(keys)

	for i := 0; i < dNum; i++ {
		tmp := make(map[string]interface{}, fieldNum)
		d := v.Index(i).Interface()
		v := reflect.ValueOf(d)
		for _, k := range keys {
			switch dstT {
			case reflect.Map:
				{
					kV := reflect.ValueOf(k)
					if v.MapIndex(kV).IsValid() {
						tmp[k] = v.MapIndex(kV).Interface()
					} else {
						tmp[k] = nil
					}
				}
			case reflect.Struct:
				{
					if v.FieldByName(k).IsValid() {
						tmp[k] = v.FieldByName(k).Interface()
					} else {
						tmp[k] = nil
					}
				}
			}
		}
		dst[i] = tmp
	}
	return dst
}

func ArrayOmit(data interface{}, keys []string) ArrayPickS {
	return ArrayPick(data, GetTypeFieldBySet(ArrayFirst(data), keys))
}

func Set() bool {
	AssetsError(errors.New("不要尝试这样, 宝贝"))
	return true
}

func Has(data interface{}, path string) bool {
	_, exist := Get(data, path)
	return exist
}

func Get(data interface{}, path string) (interface{}, bool) {
	paths := strings.SplitN(path, ".", 2)
	v := TryInterfacePtr(data)
	kind := v.Kind()
	k := paths[0]
	shouldNext := len(paths) > 1

	if !InArray([]reflect.Kind{reflect.Slice, reflect.Map, reflect.Struct}, kind) {
		return nil, false
	} else {
		item := reflect.Value{}
		switch kind {
		case reflect.Slice:
			{
				if index, err := strconv.Atoi(k); err != nil || index > v.Len()-1 || index < 0 {
					return nil, false
				} else {
					item = v.Index(index)
				}
			}
		case reflect.Map:
			{
				item = v.MapIndex(reflect.ValueOf(k))
			}
		case reflect.Struct:
			{
				item = v.FieldByName(k)
			}
		}

		if item.IsValid() {
			if shouldNext {
				return Get(item.Interface(), paths[1])
			}

			return item.Interface(), true
		} else {
			return nil, false
		}
	}
}

func AssetsReturn(assets bool, s interface{}, t interface{}) interface{} {
	if assets {
		return s
	} else {
		return t
	}
}

func AssetsError(err error) {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		buffer := new(strings.Builder)
		if ok {
			buffer.WriteString("file: ")
			buffer.WriteString(file)
			buffer.WriteString(" line: ")
			buffer.WriteString(strconv.Itoa(line))
			buffer.WriteString(" err: ")
		}
		buffer.WriteString(err.Error())
		panic(errors.New(buffer.String()))
	}
}
